package routes

import (
	"net/http"
	"time"

	ordemapp "consumo-real-server/internal/application/ordemabastecimento"
	"consumo-real-server/internal/shared/apperror"
)

type OrdemAbastecimentoHandler struct {
	service *ordemapp.Service
}

func NewOrdemAbastecimentoHandler(service *ordemapp.Service) *OrdemAbastecimentoHandler {
	return &OrdemAbastecimentoHandler{service: service}
}

type ordemRequestBody struct {
	EmpresaID            int64   `json:"empresa_id"`
	PatrimonioID         int64   `json:"patrimonio_id"`
	Numero               string  `json:"numero"`
	QuantidadeAutorizada float64 `json:"quantidade_autorizada"`
	DataValidade         *string `json:"data_validade"`
}

func (h *OrdemAbastecimentoHandler) parseDataValidade(w http.ResponseWriter, raw *string) (*time.Time, bool) {
	if raw == nil || *raw == "" {
		return nil, true
	}
	data, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		apperror.WriteError(w, apperror.Validation("parâmetro 'data_validade' inválido, use RFC3339 (ex.: 2026-01-01T00:00:00Z)", err))
		return nil, false
	}
	return &data, true
}

func (h *OrdemAbastecimentoHandler) create(w http.ResponseWriter, r *http.Request) {
	var body ordemRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	dataValidade, ok := h.parseDataValidade(w, body.DataValidade)
	if !ok {
		return
	}

	o, err := h.service.Create.Handle(r.Context(), ordemapp.CreateCommand{
		EmpresaID:            body.EmpresaID,
		PatrimonioID:         body.PatrimonioID,
		Numero:               body.Numero,
		QuantidadeAutorizada: body.QuantidadeAutorizada,
		DataValidade:         dataValidade,
		UsuarioID:            currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, o)
}

func (h *OrdemAbastecimentoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body ordemRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	dataValidade, ok := h.parseDataValidade(w, body.DataValidade)
	if !ok {
		return
	}

	o, err := h.service.Update.Handle(r.Context(), ordemapp.UpdateCommand{
		ID:                   id,
		PatrimonioID:         body.PatrimonioID,
		QuantidadeAutorizada: body.QuantidadeAutorizada,
		DataValidade:         dataValidade,
		UsuarioID:            currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, o)
}

func (h *OrdemAbastecimentoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), ordemapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrdemAbastecimentoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	o, err := h.service.Get.Handle(r.Context(), ordemapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, o)
}

func (h *OrdemAbastecimentoHandler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List.Handle(r.Context(), ordemapp.ListQuery{
		EmpresaID:    parseQueryInt64(r, "empresa_id"),
		PatrimonioID: parseQueryInt64(r, "patrimonio_id"),
		Status:       r.URL.Query().Get("status"),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

func (h *OrdemAbastecimentoHandler) autorizar(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Autorizar.Handle(r.Context(), ordemapp.AutorizarCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
