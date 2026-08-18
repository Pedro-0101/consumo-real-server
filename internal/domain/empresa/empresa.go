package empresa

import (
	"errors"
	"strings"

	shared "consumo-real-server/internal/shared"
	helper "consumo-real-server/internal/shared/helper"
)

var (
	ErrNomeObrigatorio = errors.New("nome é obrigatório")
	ErrCNPJInvalido    = errors.New("CNPJ inválido")
)

type Empresa struct {
	ID                 int64  `gorm:"primaryKey" json:"id"`
	Nome               string `gorm:"size:255;not null" json:"nome"`
	CNPJ               string `gorm:"size:20;not null;uniqueIndex" json:"cnpj"`
	Ativo              bool   `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewEmpresa(nome, cnpj string) (*Empresa, error) {

	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if !helper.IsValidCNPJ(cnpj) {
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
	if !helper.IsValidCNPJ(cnpj) {
		return ErrCNPJInvalido
	}

	e.Nome = strings.TrimSpace(nome)
	e.CNPJ = strings.TrimSpace(cnpj)
	return nil
}
