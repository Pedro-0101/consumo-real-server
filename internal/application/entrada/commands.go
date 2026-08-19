package entrada

import (
	"context"
	"errors"

	domainentrada "consumo-real-server/internal/domain/entrada"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para registrar uma entrada de combustível no reservatório.
type CreateCommand struct {
	EmpresaID      int64
	FornecedorID   int64
	ReservatorioID int64
	Quantidade     float64
	NotaFiscal     string
	UsuarioID      int64
}

type CreateHandler struct {
	repo       domainentrada.Repository
	reservRepo domainreservatorio.Repository
}

func NewCreateHandler(repo domainentrada.Repository, reservRepo domainreservatorio.Repository) *CreateHandler {
	return &CreateHandler{repo: repo, reservRepo: reservRepo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainentrada.Entrada, error) {
	reserv, err := h.reservRepo.FindByID(ctx, cmd.ReservatorioID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório", err)
	}

	e, err := domainentrada.NewEntrada(cmd.EmpresaID, cmd.FornecedorID, reserv, cmd.Quantidade, cmd.NotaFiscal)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	e.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, e); err != nil {
		return nil, apperror.Internal("falha ao criar entrada", err)
	}

	if err := h.reservRepo.Update(ctx, reserv); err != nil {
		return nil, apperror.Internal("falha ao atualizar nível do reservatório", err)
	}
	return e, nil
}

// UpdateCommand é o comando para atualizar os dados de uma entrada existente.
type UpdateCommand struct {
	ID         int64
	NotaFiscal string
	UsuarioID  int64
}

type UpdateHandler struct {
	repo domainentrada.Repository
}

func NewUpdateHandler(repo domainentrada.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainentrada.Entrada, error) {
	e, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainentrada.ErrNaoEncontrado) {
			return nil, apperror.NotFound("entrada não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar entrada", err)
	}

	e.AtualizarNotaFiscal(cmd.NotaFiscal)
	e.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, e); err != nil {
		return nil, apperror.Internal("falha ao atualizar entrada", err)
	}
	return e, nil
}

// DeleteCommand é o comando para excluir uma entrada.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainentrada.Repository
}

func NewDeleteHandler(repo domainentrada.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainentrada.ErrNaoEncontrado) {
			return apperror.NotFound("entrada não encontrada")
		}
		return apperror.Internal("falha ao buscar entrada", err)
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return apperror.Internal("falha ao excluir entrada", err)
	}
	return nil
}
