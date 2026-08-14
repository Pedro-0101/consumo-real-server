package combustivel

import (
	"context"
	"errors"

	domaincombustivel "consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo combustível.
type CreateCommand struct {
	EmpresaID  int64
	Nome       string
	Tipo       domaincombustivel.Tipo
	Unidade    domaincombustivel.Unidade
	Densidade  float64
	PrecoCusto float64
	PrecoVenda float64
	UsuarioID  int64
}

type CreateHandler struct {
	repo domaincombustivel.Repository
}

func NewCreateHandler(repo domaincombustivel.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domaincombustivel.Combustivel, error) {
	c, err := domaincombustivel.NewCombustivel(
		cmd.EmpresaID, cmd.Nome, cmd.Tipo, cmd.Unidade, cmd.Densidade, cmd.PrecoCusto, cmd.PrecoVenda,
	)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	c.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, c); err != nil {
		return nil, apperror.Internal("falha ao criar combustível", err)
	}
	return c, nil
}

// UpdateCommand é o comando para atualizar um combustível existente.
type UpdateCommand struct {
	ID         int64
	Nome       string
	Tipo       domaincombustivel.Tipo
	Unidade    domaincombustivel.Unidade
	Densidade  float64
	PrecoCusto float64
	PrecoVenda float64
	UsuarioID  int64
}

type UpdateHandler struct {
	repo domaincombustivel.Repository
}

func NewUpdateHandler(repo domaincombustivel.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domaincombustivel.Combustivel, error) {
	c, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domaincombustivel.ErrNaoEncontrado) {
			return nil, apperror.NotFound("combustível não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar combustível", err)
	}

	if err := c.Atualizar(cmd.Nome, cmd.Tipo, cmd.Unidade, cmd.Densidade); err != nil {
		return nil, apperror.FromDomain(err)
	}
	if err := c.AtualizarPrecos(cmd.PrecoCusto, cmd.PrecoVenda); err != nil {
		return nil, apperror.FromDomain(err)
	}
	c.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, c); err != nil {
		return nil, apperror.Internal("falha ao atualizar combustível", err)
	}
	return c, nil
}

// DeleteCommand é o comando para desativar um combustível.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domaincombustivel.Repository
}

func NewDeleteHandler(repo domaincombustivel.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	c, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domaincombustivel.ErrNaoEncontrado) {
			return apperror.NotFound("combustível não encontrado")
		}
		return apperror.Internal("falha ao buscar combustível", err)
	}

	c.Desativar()
	c.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, c); err != nil {
		return apperror.Internal("falha ao desativar combustível", err)
	}
	return nil
}
