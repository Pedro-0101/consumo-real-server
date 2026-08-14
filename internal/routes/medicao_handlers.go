package routes

import (
	"net/http"

	medicaoapp "consumo-real-server/internal/application/medicao"
	"consumo-real-server/internal/shared/apperror"
)

type MedicaoHandler struct {
	service *medicaoapp.Service
}

func NewMedicaoHandler(service *medicaoapp.Service) *MedicaoHandler {
	return &MedicaoHandler{service: service}
}

type medicaoRequestBody struct {
	EmpresaID      int64   `json:"empresa_id"`
	ReservatorioID int64   `json:"reservatorio_id"`
	NivelMedido    float64 `json:"nivel_medido"`
}

func (h *MedicaoHandler) create(w http.ResponseWriter, r *http.Request) {
	var body medicaoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	m, err := h.service.Create.Handle(r.Context(), medicaoapp.CreateCommand{
		EmpresaID:      body.EmpresaID,
		ReservatorioID: body.ReservatorioID,
		NivelMedido:    body.NivelMedido,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, m)
}

func (h *MedicaoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body medicaoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	m, err := h.service.Update.Handle(r.Context(), medicaoapp.UpdateCommand{
		ID:          id,
		NivelMedido: body.NivelMedido,
		UsuarioID:   currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, m)
}

func (h *MedicaoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), medicaoapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MedicaoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	m, err := h.service.Get.Handle(r.Context(), medicaoapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, m)
}

func (h *MedicaoHandler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List.Handle(r.Context(), medicaoapp.ListQuery{
		EmpresaID:      parseQueryInt64(r, "empresa_id"),
		ReservatorioID: parseQueryInt64(r, "reservatorio_id"),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
