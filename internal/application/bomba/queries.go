package bomba

import (
	"context"
	"errors"

	domainbomba "consumo-real-server/internal/domain/bomba"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma bomba pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainbomba.Repository
}

func NewGetHandler(repo domainbomba.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainbomba.Bomba, error) {
	b, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainbomba.ErrNaoEncontrado) {
			return nil, apperror.NotFound("bomba não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar bomba", err)
	}
	return b, nil
}

// ListQuery lista bombas aplicando os filtros informados.
type ListQuery struct {
	EmpresaID      int64
	LocalID        int64
	ReservatorioID int64
	Ativo          *bool
}

type ListHandler struct {
	repo domainbomba.Repository
}

func NewListHandler(repo domainbomba.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainbomba.Bomba, error) {
	list, err := h.repo.List(ctx, domainbomba.ListFilter{
		EmpresaID:      q.EmpresaID,
		LocalID:        q.LocalID,
		ReservatorioID: q.ReservatorioID,
		Ativo:          q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar bombas", err)
	}
	return list, nil
}
