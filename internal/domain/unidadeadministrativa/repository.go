package unidadeadministrativa

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("unidade administrativa não encontrada")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID int64
	Ativo     *bool
}

// Repository é o contrato de persistência do agregado UnidadeAdministrativa.
type Repository interface {
	Create(ctx context.Context, u *UnidadeAdministrativa) error
	Update(ctx context.Context, u *UnidadeAdministrativa) error
	FindByID(ctx context.Context, id int64) (*UnidadeAdministrativa, error)
	List(ctx context.Context, filter ListFilter) ([]UnidadeAdministrativa, error)
}
