package routes

import (
	"net/http"
	"strconv"

	patrimonioapp "consumo-real-server/internal/application/patrimonio"
	domainpatrimonio "consumo-real-server/internal/domain/patrimonio"
	"consumo-real-server/internal/shared/apperror"
)

type PatrimonioHandler struct {
	service *patrimonioapp.Service
}

func NewPatrimonioHandler(service *patrimonioapp.Service) *PatrimonioHandler {
	return &PatrimonioHandler{service: service}
}

type patrimonioRequestBody struct {
	EmpresaID     int64  `json:"empresa_id"`
	Nome          string `json:"nome"`
	Tipo          string `json:"tipo"`
	CodigoExterno string `json:"codigo_externo"`
	TipoMedicao   string `json:"tipo_medicao"`
}

func (h *PatrimonioHandler) create(w http.ResponseWriter, r *http.Request) {
	var body patrimonioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	p, err := h.service.Create.Handle(r.Context(), patrimonioapp.CreateCommand{
		EmpresaID:     body.EmpresaID,
		Nome:          body.Nome,
		Tipo:          body.Tipo,
		CodigoExterno: body.CodigoExterno,
		TipoMedicao:   domainpatrimonio.TipoMedicao(body.TipoMedicao),
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, p)
}

func (h *PatrimonioHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body patrimonioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	p, err := h.service.Update.Handle(r.Context(), patrimonioapp.UpdateCommand{
		ID:            id,
		Nome:          body.Nome,
		Tipo:          body.Tipo,
		CodigoExterno: body.CodigoExterno,
		TipoMedicao:   domainpatrimonio.TipoMedicao(body.TipoMedicao),
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, p)
}

func (h *PatrimonioHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), patrimonioapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PatrimonioHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	p, err := h.service.Get.Handle(r.Context(), patrimonioapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, p)
}

func (h *PatrimonioHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), patrimonioapp.ListQuery{
		EmpresaID:               parseQueryInt64(r, "empresa_id"),
		UnidadeAdministrativaID: parseQueryInt64(r, "unidade_administrativa_id"),
		Tipo:                    r.URL.Query().Get("tipo"),
		Ativo:                   ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
