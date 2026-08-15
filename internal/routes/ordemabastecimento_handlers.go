package routes

import (
	"net/http"
	"time"

	ordemapp "consumo-real-server/internal/application/ordemabastecimento"
	_ "consumo-real-server/internal/domain/ordemabastecimento"
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

// CreateOrdemAbastecimento cadastra uma nova ordem de abastecimento.
// @Summary Cadastrar Ordem de Abastecimento
// @Description Cria uma nova ordem de abastecimento no sistema.
// @Tags Ordens de Abastecimento
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ordem body ordemRequestBody true "Dados da ordem"
// @Success 201 {object} ordemabastecimento.OrdemAbastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento [post]
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

// UpdateOrdemAbastecimento atualiza os dados de uma ordem de abastecimento existente.
// @Summary Atualizar Ordem de Abastecimento
// @Description Atualiza os dados de uma ordem de abastecimento.
// @Tags Ordens de Abastecimento
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da ordem"
// @Param ordem body ordemRequestBody true "Dados atualizados da ordem"
// @Success 200 {object} ordemabastecimento.OrdemAbastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento/{id} [put]
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

// DeleteOrdemAbastecimento remove uma ordem de abastecimento do sistema.
// @Summary Excluir Ordem de Abastecimento
// @Description Remove uma ordem de abastecimento do sistema.
// @Tags Ordens de Abastecimento
// @Security BearerAuth
// @Param id path int true "ID da ordem"
// @Success 204 "Ordem excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento/{id} [delete]
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

// GetOrdemAbastecimento retorna os dados de uma ordem de abastecimento.
// @Summary Buscar Ordem de Abastecimento por ID
// @Description Retorna os dados completos de uma ordem de abastecimento.
// @Tags Ordens de Abastecimento
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da ordem"
// @Success 200 {object} ordemabastecimento.OrdemAbastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento/{id} [get]
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

// ListOrdensDeAbastecimento lista ordens de abastecimento com filtros opcionais.
// @Summary Listar Ordens de Abastecimento
// @Description Lista as ordens de abastecimento, podendo aplicar filtros.
// @Tags Ordens de Abastecimento
// @Produce json
// @Security BearerAuth
// @Param empresa_id query int false "ID da empresa"
// @Param patrimonio_id query int false "ID do patrimônio"
// @Param status query string false "Status da ordem (ABERTA, AUTORIZADA, CONCLUIDA, CANCELADA)"
// @Success 200 {array} ordemabastecimento.OrdemAbastecimento
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento [get]
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

// AutorizarOrdemAbastecimento autoriza uma ordem de abastecimento.
// @Summary Autorizar ordem de abastecimento
// @Description Autoriza uma ordem de abastecimento pendente.
// @Tags Ordens de Abastecimento
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da ordem"
// @Success 204 "Ordem autorizada"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /ordens-abastecimento/{id}/autorizar [patch]
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
