package medicao

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("medição não encontrada")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID      int64
	ReservatorioID int64
}

// Repository é o contrato de persistência do agregado Medicao.
type Repository interface {
	Create(ctx context.Context, m *Medicao) error
	Update(ctx context.Context, m *Medicao) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Medicao, error)
	List(ctx context.Context, filter ListFilter) ([]Medicao, error)
}
