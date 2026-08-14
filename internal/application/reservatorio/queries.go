package reservatorio

import (
	"context"
	"errors"

	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um reservatório pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainreservatorio.Repository
}

func NewGetHandler(repo domainreservatorio.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainreservatorio.Reservatorio, error) {
	r, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório", err)
	}
	return r, nil
}

// ListQuery lista reservatórios aplicando os filtros informados.
type ListQuery struct {
	EmpresaID     int64
	CombustivelID int64
	Ativo         *bool
}

type ListHandler struct {
	repo domainreservatorio.Repository
}

func NewListHandler(repo domainreservatorio.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainreservatorio.Reservatorio, error) {
	list, err := h.repo.List(ctx, domainreservatorio.ListFilter{
		EmpresaID:     q.EmpresaID,
		CombustivelID: q.CombustivelID,
		Ativo:         q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar reservatórios", err)
	}
	return list, nil
}
