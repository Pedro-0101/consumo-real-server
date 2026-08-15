package routes

import (
	"net/http"

	entradaapp "consumo-real-server/internal/application/entrada"
	_ "consumo-real-server/internal/domain/entrada"
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

// CreateEntrada registra uma nova entrada de combustível.
// @Summary Cadastrar Entrada
// @Description Registra uma nova entrada de combustível.
// @Tags Entradas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param entrada body entradaRequestBody true "Dados da entrada"
// @Success 201 {object} entrada.Entrada
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /entradas [post]
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

// UpdateEntrada atualiza os dados de uma entrada existente.
// @Summary Atualizar Entrada
// @Description Atualiza os dados de uma entrada.
// @Tags Entradas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da entrada"
// @Param entrada body entradaRequestBody true "Dados atualizados da entrada"
// @Success 200 {object} entrada.Entrada
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /entradas/{id} [put]
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

// DeleteEntrada remove uma entrada do sistema.
// @Summary Excluir Entrada
// @Description Remove uma entrada do sistema.
// @Tags Entradas
// @Security BearerAuth
// @Param id path int true "ID da entrada"
// @Success 204 "Entrada excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /entradas/{id} [delete]
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

// GetEntrada retorna os dados de uma entrada.
// @Summary Buscar Entrada por ID
// @Description Retorna os dados completos de uma entrada.
// @Tags Entradas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da entrada"
// @Success 200 {object} entrada.Entrada
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /entradas/{id} [get]
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

// ListEntradas lista entradas com filtros opcionais.
// @Summary Listar Entradas
// @Description Lista as entradas, podendo aplicar filtros.
// @Tags Entradas
// @Produce json
// @Security BearerAuth
// @Param empresa_id query int false "ID da empresa"
// @Param fornecedor_id query int false "ID do fornecedor"
// @Param reservatorio_id query int false "ID do reservatório"
// @Param combustivel_id query int false "ID do combustível"
// @Success 200 {array} entrada.Entrada
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /entradas [get]
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
