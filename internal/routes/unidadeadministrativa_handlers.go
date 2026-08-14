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
