package preco

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("preço não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID     int64
	CombustivelID int64
	Ativo         *bool
}

// Repository é o contrato de persistência do agregado Preco.
type Repository interface {
	Create(ctx context.Context, p *Preco) error
	Update(ctx context.Context, p *Preco) error
	FindByID(ctx context.Context, id int64) (*Preco, error)
	List(ctx context.Context, filter ListFilter) ([]Preco, error)
}
