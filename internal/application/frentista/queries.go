package frentista

import (
	"context"
	"errors"

	domainfrentista "consumo-real-server/internal/domain/frentista"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um frentista pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainfrentista.Repository
}

func NewGetHandler(repo domainfrentista.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainfrentista.Frentista, error) {
	f, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainfrentista.ErrNaoEncontrado) {
			return nil, apperror.NotFound("frentista não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar frentista", err)
	}
	return f, nil
}

// ListQuery lista frentistas aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	UsuarioID int64
	Matricula string
	Ativo     *bool
}

type ListHandler struct {
	repo domainfrentista.Repository
}

func NewListHandler(repo domainfrentista.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainfrentista.Frentista, error) {
	list, err := h.repo.List(ctx, domainfrentista.ListFilter{
		EmpresaID: q.EmpresaID,
		UsuarioID: q.UsuarioID,
		Matricula: q.Matricula,
		Ativo:     q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar frentistas", err)
	}
	return list, nil
}
