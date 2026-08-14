package patrimonio

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("patrimônio não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Tipo                    string
	Ativo                   *bool
}

// Repository é o contrato de persistência do agregado Patrimonio.
type Repository interface {
	Create(ctx context.Context, p *Patrimonio) error
	Update(ctx context.Context, p *Patrimonio) error
	FindByID(ctx context.Context, id int64) (*Patrimonio, error)
	List(ctx context.Context, filter ListFilter) ([]Patrimonio, error)
}
