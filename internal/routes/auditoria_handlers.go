package routes

import (
	"net/http"
	"strconv"
	"time"

	auditoriaapp "consumo-real-server/internal/application/auditoria"
	domainauditoria "consumo-real-server/internal/domain/auditoria"
	"consumo-real-server/internal/shared/apperror"
	"consumo-real-server/internal/shared/timeutil"
)

type AuditoriaHandler struct {
	service *auditoriaapp.Service
}

func NewAuditoriaHandler(service *auditoriaapp.Service) *AuditoriaHandler {
	return &AuditoriaHandler{service: service}
}

// ListAuditorias lista movimentações de auditoria com filtros opcionais.
// @Summary Listar auditorias
// @Description Lista as movimentações registradas no sistema. Usuários com empresa veem
// @Description apenas a própria empresa; usuários sem empresa (ADMIN_BASE) acessam todas
// @Description e devem informar o parâmetro empresa_id.
// @Tags Auditoria
// @Produce json
// @Security BearerAuth
// @Param empresa_id query int false "ID da empresa (obrigatório para ADMIN_BASE)"
// @Param entidade query string false "Nome da entidade (ex.: empresas, usuarios)"
// @Param operacao query string false "Tipo de operação (CREATE, UPDATE, DELETE)"
// @Param usuario_id query int false "ID do usuário responsável"
// @Param de query string false "Início do período (RFC3339, ex.: 2026-08-01T00:00:00-03:00)"
// @Param ate query string false "Fim do período (RFC3339, ex.: 2026-08-18T23:59:59-03:00)"
// @Param limit query int false "Limite de registros por página (padrão 100)"
// @Param offset query int false "Deslocamento para paginação"
// @Success 200 {array} domainauditoria.Auditoria
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /auditorias [get]
func (h *AuditoriaHandler) list(w http.ResponseWriter, r *http.Request) {
	empresaID := currentUserEmpresaID(r)
	if empresaID <= 0 {
		empresaID = parseQueryInt64(r, "empresa_id")
		if empresaID <= 0 {
			apperror.WriteError(w, apperror.Validation("parâmetro 'empresa_id' é obrigatório para usuários sem empresa", nil))
			return
		}
	}

	de, err := parseQueryTime(r, "de")
	if err != nil {
		apperror.WriteError(w, apperror.Validation("parâmetro 'de' inválido", err))
		return
	}
	ate, err := parseQueryTime(r, "ate")
	if err != nil {
		apperror.WriteError(w, apperror.Validation("parâmetro 'ate' inválido", err))
		return
	}

	limit := parseQueryInt(r, "limit", 100)
	offset := parseQueryInt(r, "offset", 0)

	list, err := h.service.List.Handle(r.Context(), auditoriaapp.ListQuery{
		EmpresaID: empresaID,
		Entidade:  r.URL.Query().Get("entidade"),
		Operacao:  domainauditoria.Operacao(r.URL.Query().Get("operacao")),
		UsuarioID: parseQueryInt64(r, "usuario_id"),
		De:        de,
		Ate:       ate,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

// parseQueryTime lê um parâmetro de data/hora no formato RFC3339, aceitando
// também apenas a data (2006-01-02). Retorna nil se o parâmetro não for enviado.
func parseQueryTime(r *http.Request, key string) (*time.Time, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}

	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", raw, timeutil.Now().Location()); err == nil {
		return &t, nil
	}
	return nil, apperror.Validation("formato esperado: RFC3339 (ex.: 2026-08-01T00:00:00-03:00) ou data (2026-08-01)", nil)
}

// parseQueryInt lê um parâmetro inteiro com valor padrão.
func parseQueryInt(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}