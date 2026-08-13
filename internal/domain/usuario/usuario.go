package usuario

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

type Papel string

const (
	PapelAdminBase     Papel = "ADMIN_BASE"
	PapelAdministrador Papel = "ADMINISTRADOR"
	PapelBombista      Papel = "BOMBISTA"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrEmailInvalido      = errors.New("e-mail inválido")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
	ErrPapelInvalido      = errors.New("papel inválido")
	ErrSenhaObrigatoria   = errors.New("senha é obrigatória")
)

type Usuario struct {
	ID        int64 `gorm:"primaryKey"`
	EmpresaID int64 `gorm:"index"`
	Nome      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null;uniqueIndex"`
	SenhaHash string `gorm:"size:255;not null"`
	Papel     Papel `gorm:"type:varchar(20);not null;index"`
	Ativo     bool `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewUsuario(nome, email, senhaHash string, papel Papel, empresaID int64) (*Usuario, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if !emailValido(email) {
		return nil, ErrEmailInvalido
	}
	if strings.TrimSpace(senhaHash) == "" {
		return nil, ErrSenhaObrigatoria
	}
	if !papel.isValid() {
		return nil, ErrPapelInvalido
	}
	if papel != PapelAdminBase && empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}

	return &Usuario{
		EmpresaID: empresaID,
		Nome:      strings.TrimSpace(nome),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		SenhaHash: senhaHash,
		Papel:     papel,
		Ativo:     true,
	}, nil
}

func NewAdminBase(nome, email, senhaHash string) (*Usuario, error) {
	return NewUsuario(nome, email, senhaHash, PapelAdminBase, 0)
}

func (u *Usuario) Desativar() {
	u.Ativo = false
}

func (p Papel) isValid() bool {
	switch p {
	case PapelAdminBase, PapelAdministrador, PapelBombista:
		return true
	}
	return false
}

func emailValido(email string) bool {
	email = strings.TrimSpace(email)
	return email != "" && strings.Contains(email, "@")
}
