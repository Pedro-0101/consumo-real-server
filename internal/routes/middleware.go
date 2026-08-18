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
//
// O contrato da API é enviar apenas o token JWT cru (o prefixo "Bearer " é
// adicionado pelo próprio back-end), mas o formato "Bearer <token>" também é
// aceito por compatibilidade (Swagger UI).
func bearerToken(r *http.Request) (string, bool) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", false
	}
	if parts := strings.Fields(header); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1], true
	}
	return header, true
}

// currentUserID retorna o ID do usuário autenticado, ou 0 se não autenticado.
func currentUserID(r *http.Request) int64 {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return 0
	}
	return claims.UsuarioID
}

// currentUserEmpresaID retorna o ID da empresa do usuário autenticado.
// Usuários sem empresa (ex.: ADMIN_BASE) retornam 0.
func currentUserEmpresaID(r *http.Request) int64 {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return 0
	}
	return claims.EmpresaID
}

// currentUserPapel retorna o papel do usuário autenticado, ou vazio se não autenticado.
func currentUserPapel(r *http.Request) string {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return ""
	}
	return claims.Papel
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
