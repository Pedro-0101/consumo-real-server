package routes

import (
	"net/http"
	"strconv"

	frentistaapp "consumo-real-server/internal/application/frentista"
	_ "consumo-real-server/internal/domain/frentista"
	"consumo-real-server/internal/shared/apperror"
)

type FrentistaHandler struct {
	service *frentistaapp.Service
}

func NewFrentistaHandler(service *frentistaapp.Service) *FrentistaHandler {
	return &FrentistaHandler{service: service}
}

type frentistaRequestBody struct {
	EmpresaID int64  `json:"empresa_id"`
	Nome      string `json:"nome"`
	Matricula string `json:"matricula"`
}

type vincularUsuarioRequestBody struct {
	UsuarioID int64 `json:"usuario_id"`
}

// CreateFrentista cadastra um novo frentista.
// @Summary Cadastrar frentista
// @Description Cria um novo frentista no sistema.
// @Tags Frentistas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param frentista body frentistaRequestBody true "Dados do frentista"
// @Success 201 {object} frentista.Frentista
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas [post]
func (h *FrentistaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body frentistaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Create.Handle(r.Context(), frentistaapp.CreateCommand{
		EmpresaID: body.EmpresaID,
		Nome:      body.Nome,
		Matricula: body.Matricula,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, f)
}

// UpdateFrentista atualiza os dados de um frentista existente.
// @Summary Atualizar frentista
// @Description Atualiza os dados de um frentista.
// @Tags Frentistas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do frentista"
// @Param frentista body frentistaRequestBody true "Dados atualizados do frentista"
// @Success 200 {object} frentista.Frentista
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas/{id} [put]
func (h *FrentistaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body frentistaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Update.Handle(r.Context(), frentistaapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		Matricula: body.Matricula,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

// DeleteFrentista remove um frentista do sistema.
// @Summary Excluir frentista
// @Description Remove um frentista do sistema.
// @Tags Frentistas
// @Security BearerAuth
// @Param id path int true "ID do frentista"
// @Success 204 "Frentista excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas/{id} [delete]
func (h *FrentistaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), frentistaapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetFrentista retorna os dados de um frentista.
// @Summary Buscar frentista por ID
// @Description Retorna os dados completos de um frentista.
// @Tags Frentistas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do frentista"
// @Success 200 {object} frentista.Frentista
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas/{id} [get]
func (h *FrentistaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	f, err := h.service.Get.Handle(r.Context(), frentistaapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

// ListFrentistas lista frentistas com filtros opcionais.
// @Summary Listar frentistas
// @Description Lista os frentistas, podendo aplicar filtros.
// @Tags Frentistas
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} frentista.Frentista
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas [get]
func (h *FrentistaHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), frentistaapp.ListQuery{
		EmpresaID: parseQueryInt64(r, "empresa_id"),
		UsuarioID: parseQueryInt64(r, "usuario_id"),
		Matricula: r.URL.Query().Get("matricula"),
		Ativo:     ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

// VincularUsuario vincula um usuário a um frentista.
// @Summary Vincular usuário a frentista
// @Description Vincula um usuário a um frentista.
// @Tags Frentistas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do frentista"
// @Param body body vincularUsuarioRequestBody true "ID do usuário a vincular"
// @Success 204 "Usuário vinculado com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /frentistas/{id}/usuario [patch]
func (h *FrentistaHandler) vincularUsuario(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body vincularUsuarioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	if err := h.service.VincularUsuario.Handle(r.Context(), frentistaapp.VincularUsuarioCommand{
		ID:        id,
		UsuarioID: body.UsuarioID,
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
