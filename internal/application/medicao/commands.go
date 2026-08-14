package medicao

import (
	"context"
	"errors"

	domainmedicao "consumo-real-server/internal/domain/medicao"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para registrar uma medição e corrigir o nível do reservatório.
type CreateCommand struct {
	EmpresaID      int64
	ReservatorioID int64
	NivelMedido    float64
	UsuarioID      int64
}

type CreateHandler struct {
	repo       domainmedicao.Repository
	reservRepo domainreservatorio.Repository
}

func NewCreateHandler(repo domainmedicao.Repository, reservRepo domainreservatorio.Repository) *CreateHandler {
	return &CreateHandler{repo: repo, reservRepo: reservRepo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainmedicao.Medicao, error) {
	reserv, err := h.reservRepo.FindByID(ctx, cmd.ReservatorioID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório", err)
	}

	m, err := domainmedicao.NewMedicao(cmd.EmpresaID, reserv, cmd.NivelMedido)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	m.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, m); err != nil {
		return nil, apperror.Internal("falha ao criar medição", err)
	}

	if err := h.reservRepo.Update(ctx, reserv); err != nil {
		return nil, apperror.Internal("falha ao corrigir nível do reservatório", err)
	}
	return m, nil
}

// UpdateCommand é o comando para atualizar os dados de uma medição existente.
type UpdateCommand struct {
	ID          int64
	NivelMedido float64
	UsuarioID   int64
}

type UpdateHandler struct {
	repo       domainmedicao.Repository
	reservRepo domainreservatorio.Repository
}

func NewUpdateHandler(repo domainmedicao.Repository, reservRepo domainreservatorio.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo, reservRepo: reservRepo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainmedicao.Medicao, error) {
	m, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainmedicao.ErrNaoEncontrado) {
			return nil, apperror.NotFound("medição não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar medição", err)
	}

	if cmd.NivelMedido < 0 {
		return nil, apperror.FromDomain(domainmedicao.ErrNivelMedidoInvalido)
	}

	m.NivelMedido = cmd.NivelMedido
	m.Diferenca = cmd.NivelMedido - m.NivelCalculado
	m.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, m); err != nil {
		return nil, apperror.Internal("falha ao atualizar medição", err)
	}
	return m, nil
}

// DeleteCommand é o comando para excluir uma medição.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainmedicao.Repository
}

func NewDeleteHandler(repo domainmedicao.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainmedicao.ErrNaoEncontrado) {
			return apperror.NotFound("medição não encontrada")
		}
		return apperror.Internal("falha ao buscar medição", err)
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return apperror.Internal("falha ao excluir medição", err)
	}
	return nil
}
