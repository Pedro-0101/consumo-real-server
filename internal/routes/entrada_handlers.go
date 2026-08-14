package routes

import (
	"net/http"

	entradaapp "consumo-real-server/internal/application/entrada"
	"consumo-real-server/internal/shared/apperror"
)

type EntradaHandler struct {
	service *entradaapp.Service
}

func NewEntradaHandler(service *entradaapp.Service) *EntradaHandler {
	return &EntradaHandler{service: service}
}

type entradaRequestBody struct {
	EmpresaID      int64   `json:"empresa_id"`
	FornecedorID   int64   `json:"fornecedor_id"`
	ReservatorioID int64   `json:"reservatorio_id"`
	Quantidade     float64 `json:"quantidade"`
	NotaFiscal     string  `json:"nota_fiscal"`
}

func (h *EntradaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body entradaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Create.Handle(r.Context(), entradaapp.CreateCommand{
		EmpresaID:      body.EmpresaID,
		FornecedorID:   body.FornecedorID,
		ReservatorioID: body.ReservatorioID,
		Quantidade:     body.Quantidade,
		NotaFiscal:     body.NotaFiscal,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, e)
}

func (h *EntradaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body entradaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Update.Handle(r.Context(), entradaapp.UpdateCommand{
		ID:         id,
		NotaFiscal: body.NotaFiscal,
		UsuarioID:  currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

func (h *EntradaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), entradaapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EntradaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	e, err := h.service.Get.Handle(r.Context(), entradaapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

func (h *EntradaHandler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List.Handle(r.Context(), entradaapp.ListQuery{
		EmpresaID:      parseQueryInt64(r, "empresa_id"),
		FornecedorID:   parseQueryInt64(r, "fornecedor_id"),
		ReservatorioID: parseQueryInt64(r, "reservatorio_id"),
		CombustivelID:  parseQueryInt64(r, "combustivel_id"),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
