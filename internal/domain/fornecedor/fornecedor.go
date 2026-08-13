package fornecedor

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
)

type Fornecedor struct {
	ID        int64
	EmpresaID int64
	Nome      string
	CNPJ      string
	Ativo     bool
	Audit     shared.AuditFields
}

func NewFornecedor(empresaID int64, nome, cnpj string) (*Fornecedor, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}

	return &Fornecedor{
		EmpresaID: empresaID,
		Nome:      strings.TrimSpace(nome),
		CNPJ:      strings.TrimSpace(cnpj),
		Ativo:     true,
	}, nil
}

func (f *Fornecedor) Desativar() {
	f.Ativo = false
}
