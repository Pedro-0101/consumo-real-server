package routes

import (
	"net/http"
	"strconv"

	frentistaapp "consumo-real-server/internal/application/frentista"
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

func (h *FrentistaHandler) vincularUsuario(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body struct {
		UsuarioID int64 `json:"usuario_id"`
	}
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
