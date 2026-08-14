package medicao

import (
	"context"
	"errors"

	domainmedicao "consumo-real-server/internal/domain/medicao"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma medição pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainmedicao.Repository
}

func NewGetHandler(repo domainmedicao.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainmedicao.Medicao, error) {
	m, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainmedicao.ErrNaoEncontrado) {
			return nil, apperror.NotFound("medição não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar medição", err)
	}
	return m, nil
}

// ListQuery lista medições aplicando os filtros informados.
type ListQuery struct {
	EmpresaID      int64
	ReservatorioID int64
}

type ListHandler struct {
	repo domainmedicao.Repository
}

func NewListHandler(repo domainmedicao.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainmedicao.Medicao, error) {
	list, err := h.repo.List(ctx, domainmedicao.ListFilter{
		EmpresaID:      q.EmpresaID,
		ReservatorioID: q.ReservatorioID,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar medições", err)
	}
	return list, nil
}
