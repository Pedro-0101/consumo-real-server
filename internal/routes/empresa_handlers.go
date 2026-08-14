package routes

import (
	"net/http"
	"strconv"

	empresaapp "consumo-real-server/internal/application/empresa"
	"consumo-real-server/internal/shared/apperror"
)

type EmpresaHandler struct {
	service *empresaapp.Service
}

func NewEmpresaHandler(service *empresaapp.Service) *EmpresaHandler {
	return &EmpresaHandler{service: service}
}

type empresaRequestBody struct {
	Nome string `json:"nome"`
	CNPJ string `json:"cnpj"`
}

func (h *EmpresaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body empresaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Create.Handle(r.Context(), empresaapp.CreateCommand{
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, e)
}

func (h *EmpresaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body empresaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Update.Handle(r.Context(), empresaapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

func (h *EmpresaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), empresaapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EmpresaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	e, err := h.service.Get.Handle(r.Context(), empresaapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

func (h *EmpresaHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), empresaapp.ListQuery{Ativo: ativo})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
