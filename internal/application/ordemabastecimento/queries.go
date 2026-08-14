package ordemabastecimento

import (
	"context"
	"errors"

	domainordem "consumo-real-server/internal/domain/ordemabastecimento"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma ordem de abastecimento pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainordem.Repository
}

func NewGetHandler(repo domainordem.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainordem.OrdemAbastecimento, error) {
	o, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainordem.ErrNaoEncontrado) {
			return nil, apperror.NotFound("ordem de abastecimento não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar ordem de abastecimento", err)
	}
	return o, nil
}

// ListQuery lista ordens de abastecimento aplicando os filtros informados.
type ListQuery struct {
	EmpresaID    int64
	PatrimonioID int64
	Status       string
}

type ListHandler struct {
	repo domainordem.Repository
}

func NewListHandler(repo domainordem.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainordem.OrdemAbastecimento, error) {
	list, err := h.repo.List(ctx, domainordem.ListFilter{
		EmpresaID:    q.EmpresaID,
		PatrimonioID: q.PatrimonioID,
		Status:       domainordem.Status(q.Status),
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar ordens de abastecimento", err)
	}
	return list, nil
}
