package empresa

import (
	"context"
	"errors"
	"strings"

	domainempresa "consumo-real-server/internal/domain/empresa"
	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar uma nova empresa.
type CreateCommand struct {
	Nome      string
	CNPJ      string
	UsuarioID int64
	Papel     domainusuario.Papel

	// Dados opcionais do primeiro administrador da empresa. Quando informados,
	// a empresa e o administrador são criados em uma única transação.
	AdminNome  string
	AdminEmail string
	AdminSenha string
}

type CreateHandler struct {
	repo       domainempresa.Repository
	onboarding OnboardingRepository
	hasher     domainusuario.PasswordHasher
}

func NewCreateHandler(repo domainempresa.Repository, onboarding OnboardingRepository, hasher domainusuario.PasswordHasher) *CreateHandler {
	return &CreateHandler{repo: repo, onboarding: onboarding, hasher: hasher}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainempresa.Empresa, error) {
	e, err := domainempresa.NewEmpresa(cmd.Nome, cmd.CNPJ)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	e.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if strings.TrimSpace(cmd.AdminNome) == "" &&
		strings.TrimSpace(cmd.AdminEmail) == "" &&
		strings.TrimSpace(cmd.AdminSenha) == "" {
		if err := h.repo.Create(ctx, e); err != nil {
			return nil, apperror.Internal("falha ao criar empresa", err)
		}
		return e, nil
	}

	if cmd.Papel != domainusuario.PapelAdminBase {
		return nil, apperror.Unauthorized("somente o administrador base pode criar empresa com administrador inicial")
	}

	if strings.TrimSpace(cmd.AdminNome) == "" ||
		strings.TrimSpace(cmd.AdminEmail) == "" ||
		strings.TrimSpace(cmd.AdminSenha) == "" {
		return nil, apperror.Validation("dados do administrador inicial incompletos", nil)
	}

	senhaHash, err := h.hasher.Hash(cmd.AdminSenha)
	if err != nil {
		return nil, apperror.Internal("falha ao gerar hash da senha do administrador", err)
	}

	adminNome := cmd.AdminNome
	adminEmail := cmd.AdminEmail
	if err := h.onboarding.CriarEmpresaComAdministrador(ctx, e, func(empresaID int64) (*domainusuario.Usuario, error) {
		u, err := domainusuario.NewUsuario(adminNome, adminEmail, senhaHash, domainusuario.PapelAdministrador, empresaID)
		if err != nil {
			return nil, apperror.FromDomain(err)
		}
		u.AuditFields = shared.NewAuditFields(cmd.UsuarioID)
		return u, nil
	}); err != nil {
		if ae, ok := apperror.As(err); ok {
			return nil, ae
		}
		return nil, apperror.Internal("falha ao criar empresa com administrador inicial", err)
	}
	return e, nil
}

// UpdateCommand é o comando para atualizar uma empresa existente.
type UpdateCommand struct {
	ID        int64
	Nome      string
	CNPJ      string
	UsuarioID int64
}

type UpdateHandler struct {
	repo domainempresa.Repository
}

func NewUpdateHandler(repo domainempresa.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainempresa.Empresa, error) {
	e, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainempresa.ErrNaoEncontrado) {
			return nil, apperror.NotFound("empresa não encontrada")
		}
		return nil, apperror.Internal("falha ao buscar empresa", err)
	}

	if err := e.Atualizar(cmd.Nome, cmd.CNPJ); err != nil {
		return nil, apperror.FromDomain(err)
	}
	e.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, e); err != nil {
		return nil, apperror.Internal("falha ao atualizar empresa", err)
	}
	return e, nil
}

// DeleteCommand é o comando para desativar uma empresa.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainempresa.Repository
}

func NewDeleteHandler(repo domainempresa.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	e, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainempresa.ErrNaoEncontrado) {
			return apperror.NotFound("empresa não encontrada")
		}
		return apperror.Internal("falha ao buscar empresa", err)
	}

	e.Desativar()
	e.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, e); err != nil {
		return apperror.Internal("falha ao desativar empresa", err)
	}
	return nil
}
