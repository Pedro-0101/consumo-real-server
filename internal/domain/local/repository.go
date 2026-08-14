package local

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("local não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Ativo                   *bool
}

// Repository é o contrato de persistência do agregado Local.
type Repository interface {
	Create(ctx context.Context, l *Local) error
	Update(ctx context.Context, l *Local) error
	FindByID(ctx context.Context, id int64) (*Local, error)
	List(ctx context.Context, filter ListFilter) ([]Local, error)
}
