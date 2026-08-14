package bomba

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("bomba não encontrada")

// ErrBicoNaoEncontrado é retornado quando o bico não pertence à bomba.
var ErrBicoNaoEncontrado = errors.New("bico não encontrado")

// ListFilter define os filtros disponíveis para a consulta.
type ListFilter struct {
	EmpresaID      int64
	LocalID        int64
	ReservatorioID int64
	Ativo          *bool
}

// Repository é o contrato de persistência do agregado Bomba e seus Bicos.
type Repository interface {
	Create(ctx context.Context, b *Bomba) error
	Update(ctx context.Context, b *Bomba) error
	FindByID(ctx context.Context, id int64) (*Bomba, error)
	List(ctx context.Context, filter ListFilter) ([]Bomba, error)
	AdicionarBico(ctx context.Context, bico *Bico) error
	DesativarBico(ctx context.Context, bombaID, bicoID int64) error
}
