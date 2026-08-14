package patrimonio

import (
	"context"
	"errors"

	domainpatrimonio "consumo-real-server/internal/domain/patrimonio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo patrimônio.
type CreateCommand struct {
	EmpresaID     int64
	Nome          string
	Tipo          string
	CodigoExterno string
	TipoMedicao   domainpatrimonio.TipoMedicao
	UsuarioID     int64
}

type CreateHandler struct {
	repo domainpatrimonio.Repository
}

func NewCreateHandler(repo domainpatrimonio.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainpatrimonio.Patrimonio, error) {
	p, err := domainpatrimonio.NewPatrimonio(cmd.EmpresaID, cmd.Nome, cmd.Tipo, cmd.CodigoExterno, cmd.TipoMedicao)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	p.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, p); err != nil {
		return nil, apperror.Internal("falha ao criar patrimônio", err)
	}
	return p, nil
}

// UpdateCommand é o comando para atualizar um patrimônio existente.
type UpdateCommand struct {
	ID            int64
	Nome          string
	Tipo          string
	CodigoExterno string
	TipoMedicao   domainpatrimonio.TipoMedicao
	UsuarioID     int64
}

type UpdateHandler struct {
	repo domainpatrimonio.Repository
}

func NewUpdateHandler(repo domainpatrimonio.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainpatrimonio.Patrimonio, error) {
	p, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainpatrimonio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("patrimônio não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar patrimônio", err)
	}

	if err := p.Atualizar(cmd.Nome, cmd.Tipo, cmd.CodigoExterno, cmd.TipoMedicao); err != nil {
		return nil, apperror.FromDomain(err)
	}
	p.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, p); err != nil {
		return nil, apperror.Internal("falha ao atualizar patrimônio", err)
	}
	return p, nil
}

// DeleteCommand é o comando para desativar um patrimônio.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainpatrimonio.Repository
}

func NewDeleteHandler(repo domainpatrimonio.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	p, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainpatrimonio.ErrNaoEncontrado) {
			return apperror.NotFound("patrimônio não encontrado")
		}
		return apperror.Internal("falha ao buscar patrimônio", err)
	}

	p.Desativar()
	p.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, p); err != nil {
		return apperror.Internal("falha ao desativar patrimônio", err)
	}
	return nil
}
