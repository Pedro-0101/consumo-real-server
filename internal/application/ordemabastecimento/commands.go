package ordemabastecimento

import (
	"context"
	"errors"
	"time"

	domainordem "consumo-real-server/internal/domain/ordemabastecimento"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar uma nova ordem de abastecimento.
type CreateCommand struct {
	EmpresaID            int64
	PatrimonioID         int64
	Numero               string
	QuantidadeAutorizada float64
	DataValidade         *time.Time
	UsuarioID            int64
}

type CreateHandler struct {
	repo domainordem.Repository
}

func NewCreateHandler(repo domainordem.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainordem.OrdemAbastecimento, error) {
	o, err := domainordem.NewOrdemAbastecimento(cmd.EmpresaID, cmd.PatrimonioID, cmd.Numero, cmd.QuantidadeAutorizada, cmd.DataValidade)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	o.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, o); err != nil {
		return nil, apperror.Internal("falha ao criar ordem de abastecimento", err)
	}
	return o, nil
}

// UpdateCommand é o comando para atualizar uma ordem de abastecimento existente.
type UpdateCommand struct {
	ID                   int64
	PatrimonioID         int64
	QuantidadeAutorizada float64
	DataValidade         *time.Time
	UsuarioID            int64
}

type UpdateHandler struct {
	repo domainordem.Repository
}

func NewUpdateHandler(repo domainordem.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainordem.OrdemAbastecimento, error) {
	o, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainordem.ErrNaoEncontrado) {
			return nil, apperror.NotFound("ordem de abastecimento não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar ordem de abastecimento", err)
	}

	if cmd.PatrimonioID <= 0 {
		return nil, apperror.FromDomain(domainordem.ErrPatrimonioObrigatorio)
	}
	if cmd.QuantidadeAutorizada <= 0 {
		return nil, apperror.FromDomain(domainordem.ErrQuantidadeInvalida)
	}

	o.PatrimonioID = cmd.PatrimonioID
	o.QuantidadeAutorizada = cmd.QuantidadeAutorizada
	o.DataValidade = cmd.DataValidade
	o.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, o); err != nil {
		return nil, apperror.Internal("falha ao atualizar ordem de abastecimento", err)
	}
	return o, nil
}

// DeleteCommand é o comando para cancelar uma ordem de abastecimento (soft delete).
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainordem.Repository
}

func NewDeleteHandler(repo domainordem.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	o, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainordem.ErrNaoEncontrado) {
			return apperror.NotFound("ordem de abastecimento não encontrada")
		}
		return apperror.Internal("falha ao buscar ordem de abastecimento", err)
	}

	if err := o.Cancelar(); err != nil {
		return apperror.FromDomain(err)
	}
	o.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, o); err != nil {
		return apperror.Internal("falha ao cancelar ordem de abastecimento", err)
	}
	return nil
}

// AutorizarCommand é o comando para autorizar uma ordem aberta.
type AutorizarCommand struct {
	ID        int64
	UsuarioID int64
}

type AutorizarHandler struct {
	repo domainordem.Repository
}

func NewAutorizarHandler(repo domainordem.Repository) *AutorizarHandler {
	return &AutorizarHandler{repo: repo}
}

func (h *AutorizarHandler) Handle(ctx context.Context, cmd AutorizarCommand) error {
	o, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainordem.ErrNaoEncontrado) {
			return apperror.NotFound("ordem de abastecimento não encontrada")
		}
		return apperror.Internal("falha ao buscar ordem de abastecimento", err)
	}

	if err := o.Autorizar(); err != nil {
		return apperror.FromDomain(err)
	}
	o.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, o); err != nil {
		return apperror.Internal("falha ao autorizar ordem de abastecimento", err)
	}
	return nil
}
