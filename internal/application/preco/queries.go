package preco

import (
	"context"
	"errors"

	domainpreco "consumo-real-server/internal/domain/preco"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um preço pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainpreco.Repository
}

func NewGetHandler(repo domainpreco.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainpreco.Preco, error) {
	p, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainpreco.ErrNaoEncontrado) {
			return nil, apperror.NotFound("preço não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar preço", err)
	}
	return p, nil
}

// ListQuery lista preços aplicando os filtros informados.
type ListQuery struct {
	EmpresaID     int64
	CombustivelID int64
	Ativo         *bool
}

type ListHandler struct {
	repo domainpreco.Repository
}

func NewListHandler(repo domainpreco.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainpreco.Preco, error) {
	list, err := h.repo.List(ctx, domainpreco.ListFilter{
		EmpresaID:     q.EmpresaID,
		CombustivelID: q.CombustivelID,
		Ativo:         q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar preços", err)
	}
	return list, nil
}
