package routes

import (
	"net/http"
	"strconv"

	fornecedorapp "consumo-real-server/internal/application/fornecedor"
	"consumo-real-server/internal/shared/apperror"
)

type FornecedorHandler struct {
	service *fornecedorapp.Service
}

func NewFornecedorHandler(service *fornecedorapp.Service) *FornecedorHandler {
	return &FornecedorHandler{service: service}
}

type fornecedorRequestBody struct {
	EmpresaID int64  `json:"empresa_id"`
	Nome      string `json:"nome"`
	CNPJ      string `json:"cnpj"`
}

func (h *FornecedorHandler) create(w http.ResponseWriter, r *http.Request) {
	var body fornecedorRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Create.Handle(r.Context(), fornecedorapp.CreateCommand{
		EmpresaID: body.EmpresaID,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, f)
}

func (h *FornecedorHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body fornecedorRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Update.Handle(r.Context(), fornecedorapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

func (h *FornecedorHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), fornecedorapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FornecedorHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	f, err := h.service.Get.Handle(r.Context(), fornecedorapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

func (h *FornecedorHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), fornecedorapp.ListQuery{
		EmpresaID: parseQueryInt64(r, "empresa_id"),
		CNPJ:      r.URL.Query().Get("cnpj"),
		Ativo:     ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
