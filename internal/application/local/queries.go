package local

import (
	"context"
	"errors"

	domainlocal "consumo-real-server/internal/domain/local"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um local pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainlocal.Repository
}

func NewGetHandler(repo domainlocal.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainlocal.Local, error) {
	l, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainlocal.ErrNaoEncontrado) {
			return nil, apperror.NotFound("local não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar local", err)
	}
	return l, nil
}

// ListQuery lista locais aplicando os filtros informados.
type ListQuery struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Ativo                   *bool
}

type ListHandler struct {
	repo domainlocal.Repository
}

func NewListHandler(repo domainlocal.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainlocal.Local, error) {
	list, err := h.repo.List(ctx, domainlocal.ListFilter{
		EmpresaID:               q.EmpresaID,
		UnidadeAdministrativaID: q.UnidadeAdministrativaID,
		Ativo:                   q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar locais", err)
	}
	return list, nil
}
