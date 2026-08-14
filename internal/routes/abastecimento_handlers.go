package routes

import (
	"net/http"

	abastecimentoapp "consumo-real-server/internal/application/abastecimento"
	"consumo-real-server/internal/shared/apperror"
)

type AbastecimentoHandler struct {
	service *abastecimentoapp.Service
}

func NewAbastecimentoHandler(service *abastecimentoapp.Service) *AbastecimentoHandler {
	return &AbastecimentoHandler{service: service}
}

type abastecimentoRequestBody struct {
	EmpresaID     int64   `json:"empresa_id"`
	LocalID       int64   `json:"local_id"`
	BombaID       int64   `json:"bomba_id"`
	BicoID        int64   `json:"bico_id"`
	FrentistaID   int64   `json:"frentista_id"`
	PatrimonioID  int64   `json:"patrimonio_id"`
	Quantidade    float64 `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Odometro      float64 `json:"odometro"`
	Horimetro     float64 `json:"horimetro"`
}

func (h *AbastecimentoHandler) create(w http.ResponseWriter, r *http.Request) {
	var body abastecimentoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	a, err := h.service.Create.Handle(r.Context(), abastecimentoapp.CreateCommand{
		EmpresaID:     body.EmpresaID,
		LocalID:       body.LocalID,
		BombaID:       body.BombaID,
		BicoID:        body.BicoID,
		FrentistaID:   body.FrentistaID,
		PatrimonioID:  body.PatrimonioID,
		Quantidade:    body.Quantidade,
		PrecoUnitario: body.PrecoUnitario,
		Odometro:      body.Odometro,
		Horimetro:     body.Horimetro,
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, a)
}

type transferenciaRequestBody struct {
	EmpresaID  int64   `json:"empresa_id"`
	OrigemID   int64   `json:"origem_id"`
	DestinoID  int64   `json:"destino_id"`
	Quantidade float64 `json:"quantidade"`
}

func (h *AbastecimentoHandler) createTransferencia(w http.ResponseWriter, r *http.Request) {
	var body transferenciaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	a, err := h.service.CreateTransfer.Handle(r.Context(), abastecimentoapp.CreateTransferenciaCommand{
		EmpresaID:  body.EmpresaID,
		OrigemID:   body.OrigemID,
		DestinoID:  body.DestinoID,
		Quantidade: body.Quantidade,
		UsuarioID:  currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, a)
}

type abastecimentoUpdateRequestBody struct {
	PrecoUnitario float64 `json:"preco_unitario"`
	Odometro      float64 `json:"odometro"`
	Horimetro     float64 `json:"horimetro"`
}

func (h *AbastecimentoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body abastecimentoUpdateRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	a, err := h.service.Update.Handle(r.Context(), abastecimentoapp.UpdateCommand{
		ID:            id,
		PrecoUnitario: body.PrecoUnitario,
		Odometro:      body.Odometro,
		Horimetro:     body.Horimetro,
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, a)
}

func (h *AbastecimentoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), abastecimentoapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AbastecimentoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	a, err := h.service.Get.Handle(r.Context(), abastecimentoapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, a)
}

func (h *AbastecimentoHandler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List.Handle(r.Context(), abastecimentoapp.ListQuery{
		EmpresaID:     parseQueryInt64(r, "empresa_id"),
		LocalID:       parseQueryInt64(r, "local_id"),
		BombaID:       parseQueryInt64(r, "bomba_id"),
		PatrimonioID:  parseQueryInt64(r, "patrimonio_id"),
		FrentistaID:   parseQueryInt64(r, "frentista_id"),
		CombustivelID: parseQueryInt64(r, "combustivel_id"),
		Tipo:          r.URL.Query().Get("tipo"),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
