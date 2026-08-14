package frentista

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("frentista não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID int64
	UsuarioID int64
	Matricula string
	Ativo     *bool
}

// Repository é o contrato de persistência do agregado Frentista.
type Repository interface {
	Create(ctx context.Context, f *Frentista) error
	Update(ctx context.Context, f *Frentista) error
	FindByID(ctx context.Context, id int64) (*Frentista, error)
	List(ctx context.Context, filter ListFilter) ([]Frentista, error)
}
