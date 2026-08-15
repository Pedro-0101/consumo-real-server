package routes

import (
	"net/http"

	abastecimentoapp "consumo-real-server/internal/application/abastecimento"
	_ "consumo-real-server/internal/domain/abastecimento"
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

// CreateAbastecimento registra um novo abastecimento.
// @Summary Cadastrar Abastecimento
// @Description Registra um novo abastecimento de combustível.
// @Tags Abastecimentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param abastecimento body abastecimentoRequestBody true "Dados do abastecimento"
// @Success 201 {object} abastecimento.Abastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos [post]
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

// CreateTransferencia registra uma transferência entre reservatórios.
// @Summary Registrar transferência
// @Description Registra uma transferência de combustível entre reservatórios.
// @Tags Abastecimentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transferencia body transferenciaRequestBody true "Dados da transferência"
// @Success 201 {object} abastecimento.Abastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos/transferencias [post]
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

// UpdateAbastecimento atualiza os dados de um abastecimento existente.
// @Summary Atualizar Abastecimento
// @Description Atualiza os dados de um abastecimento.
// @Tags Abastecimentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do abastecimento"
// @Param abastecimento body abastecimentoUpdateRequestBody true "Dados atualizados do abastecimento"
// @Success 200 {object} abastecimento.Abastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos/{id} [put]
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

// DeleteAbastecimento remove um abastecimento do sistema.
// @Summary Excluir Abastecimento
// @Description Remove um abastecimento do sistema.
// @Tags Abastecimentos
// @Security BearerAuth
// @Param id path int true "ID do abastecimento"
// @Success 204 "Abastecimento excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos/{id} [delete]
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

// GetAbastecimento retorna os dados de um abastecimento.
// @Summary Buscar Abastecimento por ID
// @Description Retorna os dados completos de um abastecimento.
// @Tags Abastecimentos
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do abastecimento"
// @Success 200 {object} abastecimento.Abastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos/{id} [get]
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

// ListAbastecimentos lista abastecimentos com filtros opcionais.
// @Summary Listar Abastecimentos
// @Description Lista os abastecimentos, podendo aplicar filtros.
// @Tags Abastecimentos
// @Produce json
// @Security BearerAuth
// @Param empresa_id query int false "ID da empresa"
// @Param local_id query int false "ID do local"
// @Param bomba_id query int false "ID da bomba"
// @Param patrimonio_id query int false "ID do patrimônio"
// @Param frentista_id query int false "ID do frentista"
// @Param combustivel_id query int false "ID do combustível"
// @Param tipo query string false "Tipo (ABASTECIMENTO, TRANSFERENCIA)"
// @Success 200 {array} abastecimento.Abastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /abastecimentos [get]
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
