package entrada

import (
	"context"
	"errors"

	domainentrada "consumo-real-server/internal/domain/entrada"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta uma entrada pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainentrada.Repository
}

func NewGetHandler(repo domainentrada.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainentrada.Entrada, error) {
	e, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainentrada.ErrNaoEncontrado) {
			return nil, apperror.NotFound("entrada não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar entrada", err)
	}
	return e, nil
}

// ListQuery lista entradas aplicando os filtros informados.
type ListQuery struct {
	EmpresaID      int64
	FornecedorID   int64
	ReservatorioID int64
	CombustivelID  int64
}

type ListHandler struct {
	repo domainentrada.Repository
}

func NewListHandler(repo domainentrada.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainentrada.Entrada, error) {
	list, err := h.repo.List(ctx, domainentrada.ListFilter{
		EmpresaID:      q.EmpresaID,
		FornecedorID:   q.FornecedorID,
		ReservatorioID: q.ReservatorioID,
		CombustivelID:  q.CombustivelID,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar entradas", err)
	}
	return list, nil
}
