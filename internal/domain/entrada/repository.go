package entrada

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("entrada não encontrada")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID      int64
	FornecedorID   int64
	ReservatorioID int64
	CombustivelID  int64
}

// Repository é o contrato de persistência do agregado Entrada.
type Repository interface {
	Create(ctx context.Context, e *Entrada) error
	Update(ctx context.Context, e *Entrada) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Entrada, error)
	List(ctx context.Context, filter ListFilter) ([]Entrada, error)
}
