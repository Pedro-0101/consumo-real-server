package abastecimento

import (
	"context"
	"errors"

	domainabastecimento "consumo-real-server/internal/domain/abastecimento"
	"consumo-real-server/internal/shared/apperror"
)

// GetQuery consulta um abastecimento pelo ID.
type GetQuery struct {
	ID int64
}

type GetHandler struct {
	repo domainabastecimento.Repository
}

func NewGetHandler(repo domainabastecimento.Repository) *GetHandler {
	return &GetHandler{repo: repo}
}

func (h *GetHandler) Handle(ctx context.Context, q GetQuery) (*domainabastecimento.Abastecimento, error) {
	a, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		if errors.Is(err, domainabastecimento.ErrNaoEncontrado) {
			return nil, apperror.NotFound("abastecimento não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar abastecimento", err)
	}
	return a, nil
}

// ListQuery lista abastecimentos aplicando os filtros informados.
type ListQuery struct {
	EmpresaID     int64
	LocalID       int64
	BombaID       int64
	PatrimonioID  int64
	FrentistaID   int64
	CombustivelID int64
	Tipo          string
}

type ListHandler struct {
	repo domainabastecimento.Repository
}

func NewListHandler(repo domainabastecimento.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainabastecimento.Abastecimento, error) {
	list, err := h.repo.List(ctx, domainabastecimento.ListFilter{
		EmpresaID:     q.EmpresaID,
		LocalID:       q.LocalID,
		BombaID:       q.BombaID,
		PatrimonioID:  q.PatrimonioID,
		FrentistaID:   q.FrentistaID,
		CombustivelID: q.CombustivelID,
		Tipo:          domainabastecimento.Tipo(q.Tipo),
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar abastecimentos", err)
	}
	return list, nil
}
