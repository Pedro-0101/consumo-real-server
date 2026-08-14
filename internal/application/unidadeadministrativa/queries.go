package unidadeadministrativa

import (
	"context"
	"errors"

	domainunidade "consumo-real-server/internal/domain/unidadeadministrativa"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma unidade administrativa pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainunidade.Repository
}

func NewGetHandler(repo domainunidade.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainunidade.UnidadeAdministrativa, error) {
	u, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainunidade.ErrNaoEncontrado) {
			return nil, apperror.NotFound("unidade administrativa não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar unidade administrativa", err)
	}
	return u, nil
}

// ListQuery lista unidades administrativas aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	Ativo     *bool
}

type ListHandler struct {
	repo domainunidade.Repository
}

func NewListHandler(repo domainunidade.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainunidade.UnidadeAdministrativa, error) {
	list, err := h.repo.List(ctx, domainunidade.ListFilter{
		EmpresaID: q.EmpresaID,
		Ativo:     q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar unidades administrativas", err)
	}
	return list, nil
}
