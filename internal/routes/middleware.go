package routes

import (
	"net/http"
	"strings"

	"consumo-real-server/internal/shared/apperror"
	"consumo-real-server/internal/shared/auth"
)

// middlewareAutenticacao valida o token Bearer e injeta as claims no contexto.
func middlewareAutenticacao(tokens auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				apperror.WriteError(w, apperror.Unauthorized("token de autenticação ausente"))
				return
			}

			claims, err := tokens.Validar(token)
			if err != nil {
				apperror.WriteError(w, apperror.Unauthorized("token inválido ou expirado"))
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.ContextWithClaims(r.Context(), claims)))
		})
	}
}

// bearerToken extrai o token do header Authorization (formato "Bearer <token>").
func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// currentUserID retorna o ID do usuário autenticado, ou 0 se não autenticado.
func currentUserID(r *http.Request) int64 {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return 0
	}
	return claims.UsuarioID
}
