package fornecedor

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
	helper "consumo-real-server/internal/shared/helper"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
	ErrCNPJInvalido       = errors.New("CNPJ inválido")
)

type Fornecedor struct {
	ID                 int64  `gorm:"primaryKey" json:"id"`
	EmpresaID          int64  `gorm:"not null;uniqueIndex:uk_empresa_cnpj,priority:1" json:"empresaID"`
	Nome               string `gorm:"size:255;not null" json:"nome"`
	CNPJ               *string `gorm:"size:20;uniqueIndex:uk_empresa_cnpj,priority:2" json:"cnpj"`
	Ativo              bool   `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewFornecedor(empresaID int64, nome, cnpj string) (*Fornecedor, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	cnpjNorm, err := normalizarCNPJ(cnpj)
	if err != nil {
		return nil, err
	}

	return &Fornecedor{
		EmpresaID: empresaID,
		Nome:      strings.TrimSpace(nome),
		CNPJ:      cnpjNorm,
		Ativo:     true,
	}, nil
}

func (f *Fornecedor) Desativar() {
	f.Ativo = false
}

func (f *Fornecedor) Atualizar(nome, cnpj string) error {
	if strings.TrimSpace(nome) == "" {
		return ErrNomeObrigatorio
	}
	cnpjNorm, err := normalizarCNPJ(cnpj)
	if err != nil {
		return err
	}

	f.Nome = strings.TrimSpace(nome)
	f.CNPJ = cnpjNorm
	return nil
}

// normalizarCNPJ valida o CNPJ (quando informado) e retorna apenas os dígitos.
// Valores vazios viram NULL (*string nil), permitindo múltiplos fornecedores sem CNPJ.
func normalizarCNPJ(cnpj string) (*string, error) {
	cnpj = strings.TrimSpace(cnpj)
	if cnpj == "" {
		return nil, nil
	}
	if !helper.IsValidCNPJ(cnpj) {
		return nil, ErrCNPJInvalido
	}
	var b strings.Builder
	for _, r := range cnpj {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	digits := b.String()
	return &digits, nil
}
