package fornecedor

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("fornecedor não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID int64
	CNPJ      string
	Ativo     *bool
}

// Repository é o contrato de persistência do agregado Fornecedor.
type Repository interface {
	Create(ctx context.Context, f *Fornecedor) error
	Update(ctx context.Context, f *Fornecedor) error
	FindByID(ctx context.Context, id int64) (*Fornecedor, error)
	List(ctx context.Context, filter ListFilter) ([]Fornecedor, error)
}
