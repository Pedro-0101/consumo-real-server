package local

import (
	"context"
	"errors"

	domainlocal "consumo-real-server/internal/domain/local"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo local.
type CreateCommand struct {
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Nome                    string
	Descricao               string
	UsuarioID               int64
}

type CreateHandler struct {
	repo domainlocal.Repository
}

func NewCreateHandler(repo domainlocal.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainlocal.Local, error) {
	l, err := domainlocal.NewLocal(cmd.EmpresaID, cmd.UnidadeAdministrativaID, cmd.Nome, cmd.Descricao)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	l.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, l); err != nil {
		return nil, apperror.Internal("falha ao criar local", err)
	}
	return l, nil
}

// UpdateCommand é o comando para atualizar um local existente.
type UpdateCommand struct {
	ID                      int64
	UnidadeAdministrativaID int64
	Nome                    string
	Descricao               string
	UsuarioID               int64
}

type UpdateHandler struct {
	repo domainlocal.Repository
}

func NewUpdateHandler(repo domainlocal.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainlocal.Local, error) {
	l, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainlocal.ErrNaoEncontrado) {
			return nil, apperror.NotFound("local não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar local", err)
	}

	if err := l.Atualizar(cmd.UnidadeAdministrativaID, cmd.Nome, cmd.Descricao); err != nil {
		return nil, apperror.FromDomain(err)
	}
	l.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, l); err != nil {
		return nil, apperror.Internal("falha ao atualizar local", err)
	}
	return l, nil
}

// DeleteCommand é o comando para desativar um local.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainlocal.Repository
}

func NewDeleteHandler(repo domainlocal.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	l, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainlocal.ErrNaoEncontrado) {
			return apperror.NotFound("local não encontrado")
		}
		return apperror.Internal("falha ao buscar local", err)
	}

	l.Desativar()
	l.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, l); err != nil {
		return apperror.Internal("falha ao desativar local", err)
	}
	return nil
}
