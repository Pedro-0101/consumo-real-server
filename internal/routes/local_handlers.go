package routes

import (
	"net/http"
	"strconv"

	localapp "consumo-real-server/internal/application/local"
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
