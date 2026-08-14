package usuario

import (
	"context"
	"errors"

	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um usuário pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainusuario.Repository
}

func NewGetHandler(repo domainusuario.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainusuario.Usuario, error) {
	u, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar usuário", err)
	}
	return u, nil
}

// ListQuery lista usuários aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	Papel     domainusuario.Papel
	Ativo     *bool
}

type ListHandler struct {
	repo domainusuario.Repository
}

func NewListHandler(repo domainusuario.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainusuario.Usuario, error) {
	list, err := h.repo.List(ctx, domainusuario.ListFilter{
		EmpresaID: q.EmpresaID,
		Papel:     q.Papel,
		Ativo:     q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar usuários", err)
	}
	return list, nil
}
