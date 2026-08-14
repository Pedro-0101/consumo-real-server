package preco

import (
	"context"
	"errors"
	"time"

	domainpreco "consumo-real-server/internal/domain/preco"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo preço.
type CreateCommand struct {
	EmpresaID      int64
	CombustivelID  int64
	PrecoCusto     float64
	PrecoVenda     float64
	VigenciaInicio time.Time
	UsuarioID      int64
}

type CreateHandler struct {
	repo domainpreco.Repository
}

func NewCreateHandler(repo domainpreco.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainpreco.Preco, error) {
	p, err := domainpreco.NewPreco(cmd.EmpresaID, cmd.CombustivelID, cmd.PrecoCusto, cmd.PrecoVenda, cmd.VigenciaInicio)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	p.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, p); err != nil {
		return nil, apperror.Internal("falha ao criar preço", err)
	}
	return p, nil
}

// UpdateCommand é o comando para atualizar um preço existente.
type UpdateCommand struct {
	ID             int64
	PrecoCusto     float64
	PrecoVenda     float64
	VigenciaInicio time.Time
	UsuarioID      int64
}

type UpdateHandler struct {
	repo domainpreco.Repository
}

func NewUpdateHandler(repo domainpreco.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainpreco.Preco, error) {
	p, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainpreco.ErrNaoEncontrado) {
			return nil, apperror.NotFound("preço não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar preço", err)
	}

	if cmd.PrecoCusto < 0 || cmd.PrecoVenda < 0 {
		return nil, apperror.FromDomain(domainpreco.ErrPrecoInvalido)
	}
	if cmd.VigenciaInicio.IsZero() {
		return nil, apperror.FromDomain(domainpreco.ErrVigenciaInvalida)
	}

	p.PrecoCusto = cmd.PrecoCusto
	p.PrecoVenda = cmd.PrecoVenda
	p.VigenciaInicio = cmd.VigenciaInicio
	p.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, p); err != nil {
		return nil, apperror.Internal("falha ao atualizar preço", err)
	}
	return p, nil
}

// DeleteCommand é o comando para encerrar um preço (soft delete).
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainpreco.Repository
}

func NewDeleteHandler(repo domainpreco.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	p, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainpreco.ErrNaoEncontrado) {
			return apperror.NotFound("preço não encontrado")
		}
		return apperror.Internal("falha ao buscar preço", err)
	}

	if err := p.Encerrar(time.Now()); err != nil {
		return apperror.FromDomain(err)
	}
	p.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, p); err != nil {
		return apperror.Internal("falha ao encerrar preço", err)
	}
	return nil
}
