package routes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	combustivelapp "consumo-real-server/internal/application/combustivel"
	domaincombustivel "consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/shared/apperror"
)

type CombustivelHandler struct {
	service *combustivelapp.Service
}

func NewCombustivelHandler(service *combustivelapp.Service) *CombustivelHandler {
	return &CombustivelHandler{service: service}
}

type combustivelRequestBody struct {
	EmpresaID  int64   `json:"empresa_id"`
	Nome       string  `json:"nome"`
	Tipo       string  `json:"tipo"`
	Unidade    string  `json:"unidade"`
	Densidade  float64 `json:"densidade"`
	PrecoCusto float64 `json:"preco_custo"`
	PrecoVenda float64 `json:"preco_venda"`
}

func (h *CombustivelHandler) create(w http.ResponseWriter, r *http.Request) {
	var body combustivelRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	c, err := h.service.Create.Handle(r.Context(), combustivelapp.CreateCommand{
		EmpresaID:  body.EmpresaID,
		Nome:       body.Nome,
		Tipo:       domaincombustivel.Tipo(body.Tipo),
		Unidade:    domaincombustivel.Unidade(body.Unidade),
		Densidade:  body.Densidade,
		PrecoCusto: body.PrecoCusto,
		PrecoVenda: body.PrecoVenda,
		UsuarioID:  currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, c)
}

func (h *CombustivelHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body combustivelRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	c, err := h.service.Update.Handle(r.Context(), combustivelapp.UpdateCommand{
		ID:         id,
		Nome:       body.Nome,
		Tipo:       domaincombustivel.Tipo(body.Tipo),
		Unidade:    domaincombustivel.Unidade(body.Unidade),
		Densidade:  body.Densidade,
		PrecoCusto: body.PrecoCusto,
		PrecoVenda: body.PrecoVenda,
		UsuarioID:  currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, c)
}

func (h *CombustivelHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), combustivelapp.DeleteCommand{ID: id, UsuarioID: currentUserID(r)}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CombustivelHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	c, err := h.service.Get.Handle(r.Context(), combustivelapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, c)
}

func (h *CombustivelHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), combustivelapp.ListQuery{Ativo: ativo})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

func pathID(r *http.Request) (int64, bool) {
	raw := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
