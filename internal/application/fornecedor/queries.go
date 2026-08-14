package fornecedor

import (
	"context"
	"errors"

	domainfornecedor "consumo-real-server/internal/domain/fornecedor"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um fornecedor pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainfornecedor.Repository
}

func NewGetHandler(repo domainfornecedor.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainfornecedor.Fornecedor, error) {
	f, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainfornecedor.ErrNaoEncontrado) {
			return nil, apperror.NotFound("fornecedor não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar fornecedor", err)
	}
	return f, nil
}

// ListQuery lista fornecedores aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	CNPJ      string
	Ativo     *bool
}

type ListHandler struct {
	repo domainfornecedor.Repository
}

func NewListHandler(repo domainfornecedor.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainfornecedor.Fornecedor, error) {
	list, err := h.repo.List(ctx, domainfornecedor.ListFilter{
		EmpresaID: q.EmpresaID,
		CNPJ:      q.CNPJ,
		Ativo:     q.Ativo,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar fornecedores", err)
	}
	return list, nil
}
