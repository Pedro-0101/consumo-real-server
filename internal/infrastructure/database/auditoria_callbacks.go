package database

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domainauditoria "consumo-real-server/internal/domain/auditoria"
	"consumo-real-server/internal/shared/auth"
)

// tabelaAuditoria é a tabela da própria auditoria, que fica fora do escopo
// dos callbacks para evitar recursão.
const tabelaAuditoria = "auditorias"

// chave usada para transportar o snapshot "antes" entre os callbacks
// Before/After de uma mesma operação.
const chaveSnapshotAntes = "auditoria:antes"

// RegistrarAuditoria registra os callbacks do GORM que gravam automaticamente
// a auditoria (CREATE/UPDATE/DELETE) de todas as tabelas do sistema, dentro da
// mesma transação da operação. Deve ser chamado após o AutoMigrate.
func RegistrarAuditoria(db *gorm.DB) error {
	if err := db.Callback().Create().After("gorm:create").Register("auditoria:apos_criar", aposCriar); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("auditoria:antes_atualizar", antesAtualizar); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("auditoria:apos_atualizar", aposAtualizar); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("auditoria:antes_excluir", antesExcluir); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("auditoria:apos_excluir", aposExcluir); err != nil {
		return err
	}
	return nil
}

func aposCriar(tx *gorm.DB) {
	if !auditar(tx.Statement) {
		return
	}
	novos, err := serializar(tx.Statement)
	if err != nil {
		tx.AddError(err)
		return
	}
	gravarAuditoria(tx, domainauditoria.OperacaoCreate, nil, novos)
}

func antesAtualizar(tx *gorm.DB) {
	if !auditar(tx.Statement) {
		return
	}
	antes, err := lerAntes(tx)
	if err != nil {
		tx.AddError(err)
		return
	}
	tx.Statement.Settings.Store(chaveSnapshotAntes, antes)
}

func aposAtualizar(tx *gorm.DB) {
	if !auditar(tx.Statement) {
		return
	}
	antes, _ := tx.Statement.Settings.Load(chaveSnapshotAntes)
	antesRaw, _ := antes.(json.RawMessage)

	novos, err := serializar(tx.Statement)
	if err != nil {
		tx.AddError(err)
		return
	}
	gravarAuditoria(tx, domainauditoria.OperacaoUpdate, antesRaw, novos)
}

func antesExcluir(tx *gorm.DB) {
	if !auditar(tx.Statement) {
		return
	}
	antes, err := lerAntes(tx)
	if err != nil {
		tx.AddError(err)
		return
	}
	tx.Statement.Settings.Store(chaveSnapshotAntes, antes)
}

func aposExcluir(tx *gorm.DB) {
	if !auditar(tx.Statement) {
		return
	}
	antes, _ := tx.Statement.Settings.Load(chaveSnapshotAntes)
	antesRaw, _ := antes.(json.RawMessage)
	gravarAuditoria(tx, domainauditoria.OperacaoDelete, antesRaw, nil)
}

// auditar indica se a operação na tabela informada deve ser auditada.
func auditar(stmt *gorm.Statement) bool {
	if stmt == nil || stmt.Schema == nil || stmt.Schema.Table == "" {
		return false
	}
	return stmt.Schema.Table != tabelaAuditoria
}

// gravarAuditoria monta e insere o registro de auditoria dentro da mesma
// transação da operação.
func gravarAuditoria(tx *gorm.DB, operacao domainauditoria.Operacao, antes, novos json.RawMessage) {
	stmt := tx.Statement
	a, err := domainauditoria.NovaAuditoria(
		empresaIDDoStmt(stmt),
		entidadeIDDoStmt(stmt),
		stmt.Schema.Table,
		operacao,
		antes,
		novos,
		usuarioIDDoCtx(stmt.Context),
	)
	if err != nil {
		tx.AddError(err)
		return
	}

	auditoriaTx := tx.Session(&gorm.Session{NewDB: true})
	if err := NewAuditoriaGORMRepository(auditoriaTx).Create(stmt.Context, a); err != nil {
		tx.AddError(err)
	}
}

// lerAntes consulta o estado atual (anterior) do registro no banco usando as
// condições WHERE do próprio statement (funciona para Save e Delete).
func lerAntes(tx *gorm.DB) (json.RawMessage, error) {
	stmt := tx.Statement
	if stmt.Schema == nil || stmt.Schema.Table == "" {
		return nil, nil
	}

	consulta := tx.Session(&gorm.Session{NewDB: true}).Table(stmt.Schema.Table)
	if where, ok := stmt.Clauses["WHERE"]; ok {
		if w, ok2 := where.Expression.(clause.Where); ok2 {
			for _, expr := range w.Exprs {
				consulta = consulta.Where(expr)
			}
		}
	}

	var antes map[string]any
	if err := consulta.Take(&antes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	redigirSegredos(antes)
	return json.Marshal(antes)
}

// serializar converte o modelo da operação em JSON para o snapshot "depois".
func serializar(stmt *gorm.Statement) (json.RawMessage, error) {
	rv := stmt.ReflectValue
	if !rv.IsValid() {
		return nil, nil
	}
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, nil
	}
	return json.Marshal(rv.Interface())
}

// empresaIDDoStmt extrai a empresa afetada. Entidades com campo EmpresaID usam
// o próprio campo; entidades sem o campo (ex.: empresa) usam o próprio registro.
func empresaIDDoStmt(stmt *gorm.Statement) int64 {
	rv := stmt.ReflectValue
	if rv.IsValid() {
		for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				break
			}
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			if f := rv.FieldByName("EmpresaID"); f.IsValid() && f.Kind() == reflect.Int64 {
				return f.Int()
			}
		}
	}

	if id := entidadeIDDoStmt(stmt); id > 0 {
		return id
	}
	return 0
}

// entidadeIDDoStmt retorna o ID do registro afetado pela operação.
func entidadeIDDoStmt(stmt *gorm.Statement) int64 {
	return pkValor(stmt)
}

// pkValor retorna o valor da chave primária do registro afetado, primeiro a
// partir do modelo (Save/Create) e, se necessário, das condições WHERE do
// statement (ex.: Delete(&modelo{}, id)).
func pkValor(stmt *gorm.Statement) int64 {
	if stmt.Schema == nil || stmt.Schema.PrioritizedPrimaryField == nil {
		return 0
	}
	pkName := stmt.Schema.PrioritizedPrimaryField.DBName

	rv := stmt.ReflectValue
	if rv.IsValid() {
		for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				break
			}
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			if f := rv.FieldByName(stmt.Schema.PrioritizedPrimaryField.Name); f.IsValid() && f.Kind() == reflect.Int64 {
				return f.Int()
			}
		}
	}

	if where, ok := stmt.Clauses["WHERE"]; ok {
		if w, ok2 := where.Expression.(clause.Where); ok2 {
			for _, expr := range w.Exprs {
				if eq, ok3 := expr.(clause.Eq); ok3 {
					if col, ok4 := eq.Column.(string); ok4 && col == pkName {
						if v, ok5 := eq.Value.(int64); ok5 {
							return v
						}
					}
				}
			}
		}
	}
	return 0
}

// usuarioIDDoCtx extrai o usuário autenticado do contexto da requisição.
// Operações internas (ex.: seeds) ficam com UsuarioID = 0.
func usuarioIDDoCtx(ctx context.Context) int64 {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return 0
	}
	return claims.UsuarioID
}

// redigirSegredos remove campos sensíveis (senha/hash/token) dos snapshots.
func redigirSegredos(m map[string]any) {
	for chave := range m {
		lc := strings.ToLower(chave)
		if strings.Contains(lc, "senha") ||
			strings.Contains(lc, "hash") ||
			strings.Contains(lc, "password") ||
			strings.Contains(lc, "token") ||
			strings.Contains(lc, "secret") {
			delete(m, chave)
		}
	}
}