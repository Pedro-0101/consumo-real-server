package combustivel

import (
	"context"
	"errors"

	domaincombustivel "consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um combustível pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domaincombustivel.Repository
}

func NewGetHandler(repo domaincombustivel.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domaincombustivel.Combustivel, error) {
	c, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domaincombustivel.ErrNaoEncontrado) {
			return nil, apperror.NotFound("combustível não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar combustível", err)
	}
	return c, nil
}

// ListQuery lista combustíveis aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	Ativo     *bool
}

type ListHandler struct {
	repo domaincombustivel.Repository
}

func NewListHandler(repo domaincombustivel.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domaincombustivel.Combustivel, error) {
	list, err := h.repo.List(ctx, domaincombustivel.ListFilter{
		EmpresaID: q.EmpresaID,
		Ativo:     q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar combustíveis", err)
	}
	return list, nil
}
