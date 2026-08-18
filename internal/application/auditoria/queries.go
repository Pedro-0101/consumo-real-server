package auditoria

import (
	"context"
	"time"

	domainauditoria "consumo-real-server/internal/domain/auditoria"
	"consumo-real-server/internal/shared/apperror"
)

// ListQuery lista movimentações de auditoria aplicando os filtros informados.
type ListQuery struct {
	EmpresaID int64
	Entidade  string
	Operacao  domainauditoria.Operacao
	UsuarioID int64
	De        *time.Time
	Ate       *time.Time
	Limit     int
	Offset    int
}

type ListHandler struct {
	repo domainauditoria.Repository
}

func NewListHandler(repo domainauditoria.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, q ListQuery) ([]domainauditoria.Auditoria, error) {
	list, err := h.repo.List(ctx, domainauditoria.ListFilter{
		EmpresaID: q.EmpresaID,
		Entidade:  q.Entidade,
		Operacao:  q.Operacao,
		UsuarioID: q.UsuarioID,
		De:        q.De,
		Ate:       q.Ate,
		Limit:     q.Limit,
		Offset:    q.Offset,
	})
	if err != nil {
		return nil, apperror.Internal("falha ao listar auditorias", err)
	}
	return list, nil
}