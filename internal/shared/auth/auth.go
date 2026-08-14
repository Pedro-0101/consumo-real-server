package auth

import (
	"context"
	"errors"
)

// Claims contém as informações do usuário autenticado extraídas do token.
type Claims struct {
	UsuarioID int64
	EmpresaID int64
	Papel     string
}

// ErrTokenInvalido é retornado quando o token é inválido ou expirado.
var ErrTokenInvalido = errors.New("token inválido ou expirado")

// TokenManager é o contrato de emissão e validação de tokens de autenticação.
type TokenManager interface {
	Gerar(claims Claims) (string, error)
	Validar(token string) (Claims, error)
}

// chave para armazenar as claims no contexto da requisição.
type contextKey struct{}

// ContextWithClaims insere as claims autenticadas no contexto.
func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// ClaimsFromContext extrai as claims autenticadas do contexto.
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(Claims)
	return claims, ok
}
