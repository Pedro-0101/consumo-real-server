package empresa

import (
	"context"
	"errors"

	domainempresa "consumo-real-server/internal/domain/empresa"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma empresa pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainempresa.Repository
}

func NewGetHandler(repo domainempresa.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainempresa.Empresa, error) {
	e, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainempresa.ErrNaoEncontrado) {
			return nil, apperror.NotFound("empresa não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar empresa", err)
	}
	return e, nil
}

// ListQuery lista empresas aplicando os filtros informados.
type ListQuery struct {
	Ativo *bool
}

type ListHandler struct {
	repo domainempresa.Repository
}

func NewListHandler(repo domainempresa.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainempresa.Empresa, error) {
	list, err := h.repo.List(ctx, domainempresa.ListFilter{Ativo: q.Ativo})
	if err != nil {
		return nil, apperror.Internal("falha ao listar empresas", err)
	}
	return list, nil
}
