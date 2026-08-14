package empresa

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio = errors.New("nome é obrigatório")
	ErrCNPJInvalido    = errors.New("CNPJ inválido")
)

type Empresa struct {
	ID                 int64  `gorm:"primaryKey"`
	Nome               string `gorm:"size:255;not null"`
	CNPJ               string `gorm:"size:20;not null;uniqueIndex"`
	Ativo              bool   `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewEmpresa(nome, cnpj string) (*Empresa, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if !cnpjValido(cnpj) {
		return nil, ErrCNPJInvalido
	}

	return &Empresa{
		Nome:  strings.TrimSpace(nome),
		CNPJ:  strings.TrimSpace(cnpj),
		Ativo: true,
	}, nil
}

func (e *Empresa) Desativar() {
	e.Ativo = false
}

func (e *Empresa) Atualizar(nome, cnpj string) error {
	if strings.TrimSpace(nome) == "" {
		return ErrNomeObrigatorio
	}
	if !cnpjValido(cnpj) {
		return ErrCNPJInvalido
	}

	e.Nome = strings.TrimSpace(nome)
	e.CNPJ = strings.TrimSpace(cnpj)
	return nil
}

func cnpjValido(cnpj string) bool {
	var digits strings.Builder
	for _, r := range cnpj {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	return digits.Len() == 14
}
