package auditoria

import (
	"encoding/json"
	"errors"

	"consumo-real-server/internal/shared"
)

type Operacao string

const (
	OperacaoCreate Operacao = "CREATE"
	OperacaoUpdate Operacao = "UPDATE"
	OperacaoDelete Operacao = "DELETE"
)

var (
	ErrOperacaoInvalida = errors.New("operação de auditoria inválida")
)

// Auditoria registra uma movimentação (criação, atualização ou exclusão) de
// um registro de qualquer entidade do sistema, com snapshot antes/depois.
type Auditoria struct {
	ID           int64            `gorm:"primaryKey"`
	EmpresaID    int64            `gorm:"not null;index"`
	Entidade     string           `gorm:"size:100;not null;index"`
	EntidadeID   int64            `gorm:"not null;index"`
	Operacao     Operacao         `gorm:"type:varchar(10);not null"`
	DadosAntigos json.RawMessage  `gorm:"type:jsonb"`
	DadosNovos   json.RawMessage  `gorm:"type:jsonb"`
	UsuarioID    int64            `gorm:"not null;index"`
	shared.AuditFields            `gorm:"embedded;embeddedPrefix:"`
}

func NovaAuditoria(empresaID, entidadeID int64, entidade string, operacao Operacao, dadosAntigos, dadosNovos json.RawMessage, usuarioID int64) (*Auditoria, error) {
	if !operacao.isValid() {
		return nil, ErrOperacaoInvalida
	}

	return &Auditoria{
		EmpresaID:    empresaID,
		Entidade:     entidade,
		EntidadeID:   entidadeID,
		Operacao:     operacao,
		DadosAntigos: dadosAntigos,
		DadosNovos:   dadosNovos,
		UsuarioID:    usuarioID,
		AuditFields:  shared.NewAuditFields(usuarioID),
	}, nil
}

func (o Operacao) isValid() bool {
	switch o {
	case OperacaoCreate, OperacaoUpdate, OperacaoDelete:
		return true
	}
	return false
}