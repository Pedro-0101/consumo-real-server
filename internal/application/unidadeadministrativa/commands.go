package unidadeadministrativa

import (
	"context"
	"errors"

	domainunidade "consumo-real-server/internal/domain/unidadeadministrativa"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar uma nova unidade administrativa.
type CreateCommand struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Nome                    string
	Tipo                    domainunidade.Tipo
	UsuarioID               int64
}

type CreateHandler struct {
	repo domainunidade.Repository
}

func NewCreateHandler(repo domainunidade.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainunidade.UnidadeAdministrativa, error) {
	u, err := domainunidade.NewUnidadeAdministrativa(cmd.EmpresaID, cmd.UnidadeAdministrativaID, cmd.Nome, cmd.Tipo)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	u.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, u); err != nil {
		return nil, apperror.Internal("falha ao criar unidade administrativa", err)
	}
	return u, nil
}

// UpdateCommand é o comando para atualizar uma unidade administrativa.
type UpdateCommand struct {
	ID        int64
	Nome      string
	Tipo      domainunidade.Tipo
	UsuarioID int64
}

type UpdateHandler struct {
	repo domainunidade.Repository
}

func NewUpdateHandler(repo domainunidade.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainunidade.UnidadeAdministrativa, error) {
	u, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainunidade.ErrNaoEncontrado) {
			return nil, apperror.NotFound("unidade administrativa não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar unidade administrativa", err)
	}

	if err := u.Atualizar(cmd.Nome, cmd.Tipo); err != nil {
		return nil, apperror.FromDomain(err)
	}
	u.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, u); err != nil {
		return nil, apperror.Internal("falha ao atualizar unidade administrativa", err)
	}
	return u, nil
}

// DeleteCommand é o comando para desativar uma unidade administrativa.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainunidade.Repository
}

func NewDeleteHandler(repo domainunidade.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	u, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainunidade.ErrNaoEncontrado) {
			return apperror.NotFound("unidade administrativa não encontrada")
		}
		return apperror.Internal("falha ao buscar unidade administrativa", err)
	}

	u.Desativar()
	u.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, u); err != nil {
		return apperror.Internal("falha ao desativar unidade administrativa", err)
	}
	return nil
}
