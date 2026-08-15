package routes

import (
	"net/http"
	"strconv"
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

// bearerToken extrai o token do header Authorization.
// Aceita tanto "Bearer <token>" quanto o token puro ("<token>"),
// permitindo colar apenas o token no Swagger UI.
// É tolerante a chaves/parênteses acidentais envolvendo o token.
func bearerToken(r *http.Request) (string, bool) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", false
	}

	var token string
	if parts := strings.Fields(header); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		token = parts[1]
	} else {
		token = header
	}

	token = strings.TrimSpace(token)
	token = strings.Trim(token, "{}()")
	return token, token != ""
}

// currentUserID retorna o ID do usuário autenticado, ou 0 se não autenticado.
func currentUserID(r *http.Request) int64 {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return 0
	}
	return claims.UsuarioID
}

// parseQueryInt64 retorna o valor do parâmetro de consulta como int64, ou 0.
func parseQueryInt64(r *http.Request, key string) int64 {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return 0
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
