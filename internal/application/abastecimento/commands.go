package abastecimento

import (
	"context"
	"errors"

	domainabastecimento "consumo-real-server/internal/domain/abastecimento"
	domainbomba "consumo-real-server/internal/domain/bomba"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para registrar um abastecimento.
type CreateCommand struct {
	EmpresaID     int64
	LocalID       int64
	BombaID       int64
	BicoID        int64
	FrentistaID   int64
	PatrimonioID  int64
	Quantidade    float64
	PrecoUnitario float64
	Odometro      float64
	Horimetro     float64
	UsuarioID     int64
}

type CreateHandler struct {
	repo       domainabastecimento.Repository
	bombaRepo  domainbomba.Repository
	reservRepo domainreservatorio.Repository
}

func NewCreateHandler(repo domainabastecimento.Repository, bombaRepo domainbomba.Repository, reservRepo domainreservatorio.Repository) *CreateHandler {
	return &CreateHandler{repo: repo, bombaRepo: bombaRepo, reservRepo: reservRepo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainabastecimento.Abastecimento, error) {
	bomba, err := h.bombaRepo.FindByID(ctx, cmd.BombaID)
	if err != nil {
		if errors.Is(err, domainbomba.ErrNaoEncontrado) {
			return nil, apperror.NotFound("bomba não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar bomba", err)
	}

	origem, err := h.reservRepo.FindByID(ctx, bomba.ReservatorioID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório de origem não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório de origem", err)
	}

	a, err := domainabastecimento.NewAbastecimento(
		cmd.EmpresaID, cmd.LocalID, bomba, origem,
		cmd.FrentistaID, cmd.PatrimonioID, cmd.BicoID,
		cmd.Quantidade, cmd.PrecoUnitario, cmd.Odometro, cmd.Horimetro,
	)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	a.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, a); err != nil {
		return nil, apperror.Internal("falha ao criar abastecimento", err)
	}

	if err := h.reservRepo.Update(ctx, origem); err != nil {
		return nil, apperror.Internal("falha ao atualizar nível do reservatório", err)
	}
	return a, nil
}

// CreateTransferenciaCommand é o comando para registrar uma transferência entre reservatórios.
type CreateTransferenciaCommand struct {
	EmpresaID  int64
	OrigemID   int64
	DestinoID  int64
	Quantidade float64
	UsuarioID  int64
}

type CreateTransferenciaHandler struct {
	repo       domainabastecimento.Repository
	reservRepo domainreservatorio.Repository
}

func NewCreateTransferenciaHandler(repo domainabastecimento.Repository, reservRepo domainreservatorio.Repository) *CreateTransferenciaHandler {
	return &CreateTransferenciaHandler{repo: repo, reservRepo: reservRepo}
}

func (h *CreateTransferenciaHandler) Handle(ctx context.Context, cmd CreateTransferenciaCommand) (*domainabastecimento.Abastecimento, error) {
	origem, err := h.reservRepo.FindByID(ctx, cmd.OrigemID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório de origem não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório de origem", err)
	}

	destino, err := h.reservRepo.FindByID(ctx, cmd.DestinoID)
	if err != nil {
		if errors.Is(err, domainreservatorio.ErrNaoEncontrado) {
			return nil, apperror.NotFound("reservatório de destino não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar reservatório de destino", err)
	}

	a, err := domainabastecimento.NewTransferencia(cmd.EmpresaID, origem, destino, cmd.Quantidade)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	a.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, a); err != nil {
		return nil, apperror.Internal("falha ao criar transferência", err)
	}

	if err := h.reservRepo.Update(ctx, origem); err != nil {
		return nil, apperror.Internal("falha ao atualizar nível do reservatório de origem", err)
	}
	if err := h.reservRepo.Update(ctx, destino); err != nil {
		return nil, apperror.Internal("falha ao atualizar nível do reservatório de destino", err)
	}
	return a, nil
}

// UpdateCommand é o comando para atualizar dados opcionais de um abastecimento existente.
type UpdateCommand struct {
	ID            int64
	PrecoUnitario float64
	Odometro      float64
	Horimetro     float64
	UsuarioID     int64
}

type UpdateHandler struct {
	repo domainabastecimento.Repository
}

func NewUpdateHandler(repo domainabastecimento.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainabastecimento.Abastecimento, error) {
	a, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainabastecimento.ErrNaoEncontrado) {
			return nil, apperror.NotFound("abastecimento não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar abastecimento", err)
	}

	if cmd.PrecoUnitario < 0 {
		return nil, apperror.FromDomain(domainabastecimento.ErrPrecoInvalido)
	}
	if cmd.Odometro < 0 || cmd.Horimetro < 0 {
		return nil, apperror.FromDomain(domainabastecimento.ErrLeituraInvalida)
	}

	a.PrecoUnitario = cmd.PrecoUnitario
	a.ValorTotal = cmd.PrecoUnitario * a.Quantidade
	a.Odometro = cmd.Odometro
	a.Horimetro = cmd.Horimetro
	a.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, a); err != nil {
		return nil, apperror.Internal("falha ao atualizar abastecimento", err)
	}
	return a, nil
}

// DeleteCommand é o comando para excluir um abastecimento.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainabastecimento.Repository
}

func NewDeleteHandler(repo domainabastecimento.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainabastecimento.ErrNaoEncontrado) {
			return apperror.NotFound("abastecimento não encontrado")
		}
		return apperror.Internal("falha ao buscar abastecimento", err)
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return apperror.Internal("falha ao excluir abastecimento", err)
	}
	return nil
}
