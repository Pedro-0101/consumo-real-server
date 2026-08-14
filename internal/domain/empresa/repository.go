package empresa

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("empresa não encontrada")

// ListFilter define os filtros disponíveis para a consulta de empresas.
type ListFilter struct {
	Ativo *bool
}

// Repository é o contrato de persistência do agregado Empresa.
type Repository interface {
	Create(ctx context.Context, e *Empresa) error
	Update(ctx context.Context, e *Empresa) error
	FindByID(ctx context.Context, id int64) (*Empresa, error)
	List(ctx context.Context, filter ListFilter) ([]Empresa, error)
}
