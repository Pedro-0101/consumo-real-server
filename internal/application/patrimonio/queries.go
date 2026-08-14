package patrimonio

import (
	"context"
	"errors"

	domainpatrimonio "consumo-real-server/internal/domain/patrimonio"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um patrimônio pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainpatrimonio.Repository
}

func NewGetHandler(repo domainpatrimonio.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainpatrimonio.Patrimonio, error) {
	p, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainpatrimonio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("patrimônio não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar patrimônio", err)
	}
	return p, nil
}

// ListQuery lista patrimônios aplicando os filtros informados.
type ListQuery struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Tipo                    string
	Ativo                   *bool
}

type ListHandler struct {
	repo domainpatrimonio.Repository
}

func NewListHandler(repo domainpatrimonio.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainpatrimonio.Patrimonio, error) {
	list, err := h.repo.List(ctx, domainpatrimonio.ListFilter{
		EmpresaID:               q.EmpresaID,
		UnidadeAdministrativaID: q.UnidadeAdministrativaID,
		Tipo:                    q.Tipo,
		Ativo:                   q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar patrimônios", err)
	}
	return list, nil
}
