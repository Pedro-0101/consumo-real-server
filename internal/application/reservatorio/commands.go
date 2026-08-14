package reservatorio

import (
	"context"
	"errors"

	domaincombustivel "consumo-real-server/internal/domain/combustivel"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo reservatório.
type CreateCommand struct {
	EmpresaID     int64
	Nome          string
	Capacidade    float64
	NivelInicial  float64
	NivelMinimo   float64
	CombustivelID int64
	UsuarioID     int64
}

type CreateHandler struct {
	repo domainreservatorio.Repository
}

func NewCreateHandler(repo domainreservatorio.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainreservatorio.Reservatorio, error) {
	combustivel := domaincombustivel.Combustivel{ID: cmd.CombustivelID, EmpresaID: cmd.EmpresaID}
	r, err := domainreservatorio.NewReservatorio(cmd.EmpresaID, cmd.Nome, cmd.Capacidade, cmd.NivelInicial, cmd.NivelMinimo, combustivel)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	r.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, r); err != nil {
		return nil, apperror.Internal("falha ao criar reservatório", err)
	}
	return r, nil
}

// UpdateCommand é o comando para atualizar um reservatório existente.
type UpdateCommand struct {
	ID            int64
	Nome          string
	Capacidade    float64
	NivelMinimo   float64
	CombustivelID int64
	UsuarioID     int64
}

type UpdateHandler struct {
	repo domainreservatorio.Repository
}

func NewUpdateHandler(repo domainreservatorio.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainreservatorio.Reservatorio, error) {
	r, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório", err)
	}

	combustivel := domaincombustivel.Combustivel{ID: cmd.CombustivelID, EmpresaID: r.EmpresaID}
	if err := r.Atualizar(cmd.Nome, cmd.Capacidade, cmd.NivelMinimo, combustivel); err != nil {
		return nil, apperror.FromDomain(err)
	}
	r.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, r); err != nil {
		return nil, apperror.Internal("falha ao atualizar reservatório", err)
	}
	return r, nil
}

// DeleteCommand é o comando para desativar um reservatório.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainreservatorio.Repository
}

func NewDeleteHandler(repo domainreservatorio.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	r, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return apperror.NotFound("reservatório não encontrado")
		}
		return apperror.Internal("falha ao buscar reservatório", err)
	}

	r.Desativar()
	r.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, r); err != nil {
		return apperror.Internal("falha ao desativar reservatório", err)
	}
	return nil
}
