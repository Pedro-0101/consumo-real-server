package combustivel

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("combustível não encontrado")

// ListFilter define os filtros disponíveis para a consulta de combustíveis.
type ListFilter struct {
	EmpresaID int64
	Ativo     *bool
}

// Repository é o contrato de persistência do agregado Combustivel.
type Repository interface {
	Create(ctx context.Context, c *Combustivel) error
	Update(ctx context.Context, c *Combustivel) error
	FindByID(ctx context.Context, id int64) (*Combustivel, error)
	List(ctx context.Context, filter ListFilter) ([]Combustivel, error)
}
