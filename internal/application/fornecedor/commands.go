package fornecedor

import (
	"context"
	"errors"

	domainfornecedor "consumo-real-server/internal/domain/fornecedor"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo fornecedor.
type CreateCommand struct {
	EmpresaID int64
	Nome      string
	CNPJ      string
	UsuarioID int64
}

type CreateHandler struct {
	repo domainfornecedor.Repository
}

func NewCreateHandler(repo domainfornecedor.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainfornecedor.Fornecedor, error) {
	f, err := domainfornecedor.NewFornecedor(cmd.EmpresaID, cmd.Nome, cmd.CNPJ)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	f.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, f); err != nil {
		return nil, apperror.Internal("falha ao criar fornecedor", err)
	}
	return f, nil
}

// UpdateCommand é o comando para atualizar um fornecedor existente.
type UpdateCommand struct {
	ID        int64
	Nome      string
	CNPJ      string
	UsuarioID int64
}

type UpdateHandler struct {
	repo domainfornecedor.Repository
}

func NewUpdateHandler(repo domainfornecedor.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainfornecedor.Fornecedor, error) {
	f, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainfornecedor.ErrNaoEncontrado) {
			return nil, apperror.NotFound("fornecedor não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar fornecedor", err)
	}

	if err := f.Atualizar(cmd.Nome, cmd.CNPJ); err != nil {
		return nil, apperror.FromDomain(err)
	}
	f.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, f); err != nil {
		return nil, apperror.Internal("falha ao atualizar fornecedor", err)
	}
	return f, nil
}

// DeleteCommand é o comando para desativar um fornecedor.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainfornecedor.Repository
}

func NewDeleteHandler(repo domainfornecedor.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	f, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainfornecedor.ErrNaoEncontrado) {
			return apperror.NotFound("fornecedor não encontrado")
		}
		return apperror.Internal("falha ao buscar fornecedor", err)
	}

	f.Desativar()
	f.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, f); err != nil {
		return apperror.Internal("falha ao desativar fornecedor", err)
	}
	return nil
}
