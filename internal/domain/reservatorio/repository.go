package reservatorio

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("reservatório não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID     int64
	CombustivelID int64
	Ativo         *bool
}

// Repository é o contrato de persistência do agregado Reservatorio.
type Repository interface {
	Create(ctx context.Context, r *Reservatorio) error
	Update(ctx context.Context, r *Reservatorio) error
	FindByID(ctx context.Context, id int64) (*Reservatorio, error)
	List(ctx context.Context, filter ListFilter) ([]Reservatorio, error)
}
