package frentista

import (
	"context"
	"errors"

	domainfrentista "consumo-real-server/internal/domain/frentista"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo frentista.
type CreateCommand struct {
	EmpresaID int64
	Nome      string
	Matricula string
	UsuarioID int64
}

type CreateHandler struct {
	repo domainfrentista.Repository
}

func NewCreateHandler(repo domainfrentista.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainfrentista.Frentista, error) {
	f, err := domainfrentista.NewFrentista(cmd.EmpresaID, cmd.Nome, cmd.Matricula)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	f.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, f); err != nil {
		return nil, apperror.Internal("falha ao criar frentista", err)
	}
	return f, nil
}

// UpdateCommand é o comando para atualizar um frentista existente.
type UpdateCommand struct {
	ID        int64
	Nome      string
	Matricula string
	UsuarioID int64
}

type UpdateHandler struct {
	repo domainfrentista.Repository
}

func NewUpdateHandler(repo domainfrentista.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainfrentista.Frentista, error) {
	f, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainfrentista.ErrNaoEncontrado) {
			return nil, apperror.NotFound("frentista não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar frentista", err)
	}

	if err := f.Atualizar(cmd.Nome, cmd.Matricula); err != nil {
		return nil, apperror.FromDomain(err)
	}
	f.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, f); err != nil {
		return nil, apperror.Internal("falha ao atualizar frentista", err)
	}
	return f, nil
}

// DeleteCommand é o comando para desativar um frentista.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainfrentista.Repository
}

func NewDeleteHandler(repo domainfrentista.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	f, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainfrentista.ErrNaoEncontrado) {
			return apperror.NotFound("frentista não encontrado")
		}
		return apperror.Internal("falha ao buscar frentista", err)
	}

	f.Desativar()
	f.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, f); err != nil {
		return apperror.Internal("falha ao desativar frentista", err)
	}
	return nil
}

// VincularUsuarioCommand é o comando para vincular um usuário ao frentista.
type VincularUsuarioCommand struct {
	ID        int64
	UsuarioID int64
}

type VincularUsuarioHandler struct {
	repo domainfrentista.Repository
}

func NewVincularUsuarioHandler(repo domainfrentista.Repository) *VincularUsuarioHandler {
	return &VincularUsuarioHandler{repo: repo}
}

func (h *VincularUsuarioHandler) Handle(ctx context.Context, cmd VincularUsuarioCommand) error {
	f, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainfrentista.ErrNaoEncontrado) {
			return apperror.NotFound("frentista não encontrado")
		}
		return apperror.Internal("falha ao buscar frentista", err)
	}

	f.VincularUsuario(cmd.UsuarioID)

	if err := h.repo.Update(ctx, f); err != nil {
		return apperror.Internal("falha ao vincular usuário", err)
	}
	return nil
}
