package usuario

import (
	"context"
	"errors"
	"strings"

	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared/apperror"
	"consumo-real-server/internal/shared/auth"
)

// LoginCommand é o comando de autenticação de um usuário.
type LoginCommand struct {
	Email string
	Senha string
}

// LoginResult é o retorno do login com o token e o usuário autenticado.
type LoginResult struct {
	Token   string                 `json:"token"`
	Usuario *domainusuario.Usuario `json:"usuario"`
}

type LoginHandler struct {
	repo   domainusuario.Repository
	hasher domainusuario.PasswordHasher
	tokens auth.TokenManager
}

func NewLoginHandler(repo domainusuario.Repository, hasher domainusuario.PasswordHasher, tokens auth.TokenManager) *LoginHandler {
	return &LoginHandler{repo: repo, hasher: hasher, tokens: tokens}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	u, err := h.repo.FindByEmail(ctx, strings.TrimSpace(cmd.Email))
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return nil, apperror.Unauthorized("credenciais inválidas")
		}
		return nil, apperror.Internal("falha ao buscar usuário", err)
	}

	if !u.Ativo || !h.hasher.Verificar(u.SenhaHash, cmd.Senha) {
		return nil, apperror.Unauthorized("credenciais inválidas")
	}

	token, err := h.tokens.Gerar(auth.Claims{
		UsuarioID: u.ID,
		EmpresaID: u.EmpresaID,
		Papel:     string(u.Papel),
	})
	if err != nil {
		return nil, apperror.Internal("falha ao gerar token de autenticação", err)
	}

	return &LoginResult{Token: token, Usuario: u}, nil
}

// MeQuery retorna o usuário autenticado com base nas claims do token.
type MeQuery struct {
	UsuarioID int64
}

type MeHandler struct {
	repo domainusuario.Repository
}

func NewMeHandler(repo domainusuario.Repository) *MeHandler {
	return &MeHandler{repo: repo}
}

func (h *MeHandler) Handle(ctx context.Context, q MeQuery) (*domainusuario.Usuario, error) {
	u, err := h.repo.FindByID(ctx, q.UsuarioID)
	if err != nil {
		if errors.Is(err, domainusuario.ErrNaoEncontrado) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao buscar usuário", err)
	}
	return u, nil
}
