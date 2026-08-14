package usuario

import (
	"context"
	"errors"

	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/apperror"
)

// CreateCommand é o comando para cadastrar um novo usuário.
type CreateCommand struct {
	EmpresaID int64
	Nome      string
	Email     string
	Senha     string
	Papel     domainusuario.Papel
	UsuarioID int64
}

type CreateHandler struct {
	repo   domainusuario.Repository
	hasher domainusuario.PasswordHasher
}

func NewCreateHandler(repo domainusuario.Repository, hasher domainusuario.PasswordHasher) *CreateHandler {
	return &CreateHandler{repo: repo, hasher: hasher}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*domainusuario.Usuario, error) {
	senhaHash, err := h.hasher.Hash(cmd.Senha)
	if err != nil {
		return nil, apperror.Internal("falha ao gerar hash da senha", err)
	}

	u, err := domainusuario.NewUsuario(cmd.Nome, cmd.Email, senhaHash, cmd.Papel, cmd.EmpresaID)
	if err != nil {
		return nil, apperror.FromDomain(err)
	}
	u.AuditFields = shared.NewAuditFields(cmd.UsuarioID)

	if err := h.repo.Create(ctx, u); err != nil {
		return nil, apperror.Internal("falha ao criar usuário", err)
	}
	return u, nil
}

// UpdateCommand é o comando para atualizar um usuário existente.
type UpdateCommand struct {
	ID        int64
	EmpresaID int64
	Nome      string
	Email     string
	Papel     domainusuario.Papel
	UsuarioID int64
}

type UpdateHandler struct {
	repo domainusuario.Repository
}

func NewUpdateHandler(repo domainusuario.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*domainusuario.Usuario, error) {
	u, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar usuário", err)
	}

	if err := u.Atualizar(cmd.Nome, cmd.Email, cmd.Papel, cmd.EmpresaID); err != nil {
		return nil, apperror.FromDomain(err)
	}
	u.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, u); err != nil {
		return nil, apperror.Internal("falha ao atualizar usuário", err)
	}
	return u, nil
}

// ChangePasswordCommand é o comando para alterar a senha de um usuário.
type ChangePasswordCommand struct {
	ID         int64
	SenhaAtual string
	NovaSenha  string
	UsuarioID  int64
}

type ChangePasswordHandler struct {
	repo   domainusuario.Repository
	hasher domainusuario.PasswordHasher
}

func NewChangePasswordHandler(repo domainusuario.Repository, hasher domainusuario.PasswordHasher) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo, hasher: hasher}
}

func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	u, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return apperror.NotFound("usuário não encontrado")
		}
		return apperror.Internal("falha ao buscar usuário", err)
	}

	if !h.hasher.Verificar(u.SenhaHash, cmd.SenhaAtual) {
		return apperror.Unauthorized("senha atual incorreta")
	}

	novaHash, err := h.hasher.Hash(cmd.NovaSenha)
	if err != nil {
		return apperror.Internal("falha ao gerar hash da nova senha", err)
	}
	if err := u.AlterarSenha(novaHash); err != nil {
		return apperror.FromDomain(err)
	}
	u.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, u); err != nil {
		return apperror.Internal("falha ao atualizar senha", err)
	}
	return nil
}

// DeleteCommand é o comando para desativar um usuário.
type DeleteCommand struct {
	ID        int64
	UsuarioID int64
}

type DeleteHandler struct {
	repo domainusuario.Repository
}

func NewDeleteHandler(repo domainusuario.Repository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	u, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return apperror.NotFound("usuário não encontrado")
		}
		return apperror.Internal("falha ao buscar usuário", err)
	}

	u.Desativar()
	u.AuditFields.Update(cmd.UsuarioID)

	if err := h.repo.Update(ctx, u); err != nil {
		return apperror.Internal("falha ao desativar usuário", err)
	}
	return nil
}
