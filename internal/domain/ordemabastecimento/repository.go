package ordemabastecimento

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("ordem de abastecimento não encontrada")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID    int64
	PatrimonioID int64
	Status       Status
}

// Repository é o contrato de persistência do agregado OrdemAbastecimento.
type Repository interface {
	Create(ctx context.Context, o *OrdemAbastecimento) error
	Update(ctx context.Context, o *OrdemAbastecimento) error
	FindByID(ctx context.Context, id int64) (*OrdemAbastecimento, error)
	List(ctx context.Context, filter ListFilter) ([]OrdemAbastecimento, error)
}
