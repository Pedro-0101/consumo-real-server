package abastecimento

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("abastecimento não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID     int64
	LocalID       int64
	BombaID       int64
	PatrimonioID  int64
	FrentistaID   int64
	CombustivelID int64
	Tipo          Tipo
}

// Repository é o contrato de persistência do agregado Abastecimento.
type Repository interface {
	Create(ctx context.Context, a *Abastecimento) error
	Update(ctx context.Context, a *Abastecimento) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Abastecimento, error)
	List(ctx context.Context, filter ListFilter) ([]Abastecimento, error)
}
