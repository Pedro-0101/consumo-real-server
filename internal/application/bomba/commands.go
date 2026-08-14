package bomba

import (
	"context"
	"errors"

	domainbomba "consumo-real-server/internal/domain/bomba"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar uma nova bomba.
type CreateCommand struct {
	EmpresaID      int64
	LocalID        int64
	ReservatorioID int64
	Movel          bool
	Nome           string
	Descricao      string
	UsuarioID      int64
}

type CreateHandler struct {
	repo domainbomba.Repository
}

func NewCreateHandler(repo domainbomba.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainbomba.Bomba, error) {
	b, err := domainbomba.NewBomba(cmd.EmpresaID, cmd.LocalID, cmd.ReservatorioID, cmd.Movel, cmd.Nome, cmd.Descricao)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	b.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, b); err != nil {
		return nil, apperror.Internal("falha ao criar bomba", err)
	}
	return b, nil
}

// UpdateCommand é o comando para atualizar uma bomba existente.
type UpdateCommand struct {
	ID             int64
	LocalID        int64
	ReservatorioID int64
	Movel          bool
	Nome           string
	Descricao      string
	UsuarioID      int64
}

type UpdateHandler struct {
	repo domainbomba.Repository
}

func NewUpdateHandler(repo domainbomba.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainbomba.Bomba, error) {
	b, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainbomba.ErrNaoEncontrado) {
			return nil, apperror.NotFound("bomba não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar bomba", err)
	}

	if err := b.Atualizar(cmd.LocalID, cmd.ReservatorioID, cmd.Movel, cmd.Nome, cmd.Descricao); err != nil {
		return nil, apperror.FromDomain(err)
	}
	b.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, b); err != nil {
		return nil, apperror.Internal("falha ao atualizar bomba", err)
	}
	return b, nil
}

// DeleteCommand é o comando para desativar uma bomba.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainbomba.Repository
}

func NewDeleteHandler(repo domainbomba.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	b, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainbomba.ErrNaoEncontrado) {
			return apperror.NotFound("bomba não encontrada")
		}
		return apperror.Internal("falha ao buscar bomba", err)
	}

	b.Desativar()
	b.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, b); err != nil {
		return apperror.Internal("falha ao desativar bomba", err)
	}
	return nil
}

// AdicionarBicoCommand é o comando para adicionar um bico a uma bomba.
type AdicionarBicoCommand struct {
	BombaID   int64
	Nome      string
	UsuarioID int64
}

type AdicionarBicoHandler struct {
	repo domainbomba.Repository
}

func NewAdicionarBicoHandler(repo domainbomba.Repository) *AdicionarBicoHandler {
	return &AdicionarBicoHandler{repo: repo}
}

func (h *AdicionarBicoHandler) Handle(ctx context.Context, cmd AdicionarBicoCommand) (*domainbomba.Bico, error) {
	b, err := h.repo.FindByID(ctx, cmd.BombaID)
	if err != nil {
		if errors.Is(err, domainbomba.ErrNaoEncontrado) {
			return nil, apperror.NotFound("bomba não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar bomba", err)
	}

	bico, err := b.AdicionarBico(cmd.Nome)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	bico.ID = 0
	bico.BombaID = b.ID

	if err := h.repo.AdicionarBico(ctx, bico); err != nil {
		return nil, apperror.Internal("falha ao adicionar bico", err)
	}
	return bico, nil
}

// DesativarBicoCommand é o comando para desativar um bico.
type DesativarBicoCommand struct {
	BombaID   int64
	BicoID    int64
	UsuarioID int64
}

type DesativarBicoHandler struct {
	repo domainbomba.Repository
}

func NewDesativarBicoHandler(repo domainbomba.Repository) *DesativarBicoHandler {
	return &DesativarBicoHandler{repo: repo}
}

func (h *DesativarBicoHandler) Handle(ctx context.Context, cmd DesativarBicoCommand) error {
	if err := h.repo.DesativarBico(ctx, cmd.BombaID, cmd.BicoID); err != nil {
		if errors.Is(err, domainbomba.ErrBicoNaoEncontrado) {
			return apperror.NotFound("bico não encontrado")
		}
		return apperror.Internal("falha ao desativar bico", err)
	}
	return nil
}
