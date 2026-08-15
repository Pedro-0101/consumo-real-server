package routes

import (
	"net/http"
	"strconv"

	unidadeapp "consumo-real-server/internal/application/unidadeadministrativa"
	domainunidade "consumo-real-server/internal/domain/unidadeadministrativa"
	"consumo-real-server/internal/shared/apperror"
)

type UnidadeAdministrativaHandler struct {
	service *unidadeapp.Service
}

func NewUnidadeAdministrativaHandler(service *unidadeapp.Service) *UnidadeAdministrativaHandler {
	return &UnidadeAdministrativaHandler{service: service}
}

type unidadeAdministrativaRequestBody struct {
	EmpresaID               int64  `json:"empresa_id"`
	UnidadeAdministrativaID int64  `json:"unidade_administrativa_id"`
	Nome                    string `json:"nome"`
	Tipo                    string `json:"tipo"`
}

// CreateUnidadeAdministrativa cadastra uma nova unidade administrativa.
// @Summary Cadastrar Unidade Administrativa
// @Description Cria uma nova unidade administrativa no sistema.
// @Tags Unidades Administrativas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unidadeAdministrativa body unidadeAdministrativaRequestBody true "Dados da unidade administrativa"
// @Success 201 {object} domainunidade.UnidadeAdministrativa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /unidades-administrativas [post]
func (h *UnidadeAdministrativaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body unidadeAdministrativaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	u, err := h.service.Create.Handle(r.Context(), unidadeapp.CreateCommand{
		EmpresaID:               body.EmpresaID,
		UnidadeAdministrativaID: body.UnidadeAdministrativaID,
		Nome:                    body.Nome,
		Tipo:                    domainunidade.Tipo(body.Tipo),
		UsuarioID:               currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, u)
}

// UpdateUnidadeAdministrativa atualiza os dados de uma unidade administrativa existente.
// @Summary Atualizar Unidade Administrativa
// @Description Atualiza os dados de uma unidade administrativa.
// @Tags Unidades Administrativas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da unidade administrativa"
// @Param unidadeAdministrativa body unidadeAdministrativaRequestBody true "Dados atualizados da unidade administrativa"
// @Success 200 {object} domainunidade.UnidadeAdministrativa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /unidades-administrativas/{id} [put]
func (h *UnidadeAdministrativaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body unidadeAdministrativaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	u, err := h.service.Update.Handle(r.Context(), unidadeapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		Tipo:      domainunidade.Tipo(body.Tipo),
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}

// DeleteUnidadeAdministrativa remove uma unidade administrativa do sistema.
// @Summary Excluir Unidade Administrativa
// @Description Remove uma unidade administrativa do sistema.
// @Tags Unidades Administrativas
// @Security BearerAuth
// @Param id path int true "ID da unidade administrativa"
// @Success 204 "Unidade administrativa excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /unidades-administrativas/{id} [delete]
func (h *UnidadeAdministrativaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), unidadeapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetUnidadeAdministrativa retorna os dados de uma unidade administrativa.
// @Summary Buscar Unidade Administrativa por ID
// @Description Retorna os dados completos de uma unidade administrativa.
// @Tags Unidades Administrativas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da unidade administrativa"
// @Success 200 {object} domainunidade.UnidadeAdministrativa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /unidades-administrativas/{id} [get]
func (h *UnidadeAdministrativaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	u, err := h.service.Get.Handle(r.Context(), unidadeapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}

// ListUnidadesAdministrativas lista unidades administrativas com filtros opcionais.
// @Summary Listar Unidades Administrativas
// @Description Lista as unidades administrativas, podendo aplicar filtros.
// @Tags Unidades Administrativas
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} domainunidade.UnidadeAdministrativa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /unidades-administrativas [get]
func (h *UnidadeAdministrativaHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), unidadeapp.ListQuery{
		EmpresaID: parseQueryInt64(r, "empresa_id"),
		Ativo:     ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
