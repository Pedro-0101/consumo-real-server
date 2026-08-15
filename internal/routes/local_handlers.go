package routes

import (
	"net/http"
	"strconv"

	localapp "consumo-real-server/internal/application/local"
	_ "consumo-real-server/internal/domain/local"
	"consumo-real-server/internal/shared/apperror"
)

type LocalHandler struct {
	service *localapp.Service
}

func NewLocalHandler(service *localapp.Service) *LocalHandler {
	return &LocalHandler{service: service}
}

type localRequestBody struct {
	EmpresaID               int64  `json:"empresa_id"`
	UnidadeAdministrativaID int64  `json:"unidade_administrativa_id"`
	Nome                    string `json:"nome"`
	Descricao               string `json:"descricao"`
}

// CreateLocal cadastra um novo local.
// @Summary Cadastrar Local
// @Description Cria um novo local no sistema.
// @Tags Locais
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param local body localRequestBody true "Dados do local"
// @Success 201 {object} local.Local
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /locais [post]
func (h *LocalHandler) create(w http.ResponseWriter, r *http.Request) {
	var body localRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	l, err := h.service.Create.Handle(r.Context(), localapp.CreateCommand{
		EmpresaID:               body.EmpresaID,
		UnidadeAdministrativaID: body.UnidadeAdministrativaID,
		Nome:                    body.Nome,
		Descricao:               body.Descricao,
		UsuarioID:               currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, l)
}

// UpdateLocal atualiza os dados de um local existente.
// @Summary Atualizar Local
// @Description Atualiza os dados de um local.
// @Tags Locais
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do local"
// @Param local body localRequestBody true "Dados atualizados do local"
// @Success 200 {object} local.Local
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /locais/{id} [put]
func (h *LocalHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body localRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	l, err := h.service.Update.Handle(r.Context(), localapp.UpdateCommand{
		ID:                      id,
		UnidadeAdministrativaID: body.UnidadeAdministrativaID,
		Nome:                    body.Nome,
		Descricao:               body.Descricao,
		UsuarioID:               currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, l)
}

// DeleteLocal remove um local do sistema.
// @Summary Excluir Local
// @Description Remove um local do sistema.
// @Tags Locais
// @Security BearerAuth
// @Param id path int true "ID do local"
// @Success 204 "Local excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /locais/{id} [delete]
func (h *LocalHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), localapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetLocal retorna os dados de um local.
// @Summary Buscar Local por ID
// @Description Retorna os dados completos de um local.
// @Tags Locais
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do local"
// @Success 200 {object} local.Local
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /locais/{id} [get]
func (h *LocalHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	l, err := h.service.Get.Handle(r.Context(), localapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, l)
}

// ListLocais lista locais com filtros opcionais.
// @Summary Listar Locais
// @Description Lista os locais, podendo aplicar filtros.
// @Tags Locais
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} local.Local
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /locais [get]
func (h *LocalHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), localapp.ListQuery{
		EmpresaID:               parseQueryInt64(r, "empresa_id"),
		UnidadeAdministrativaID: parseQueryInt64(r, "unidade_administrativa_id"),
		Ativo:                   ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
