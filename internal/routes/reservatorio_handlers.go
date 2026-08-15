package routes

import (
	"net/http"
	"strconv"

	reservatorioapp "consumo-real-server/internal/application/reservatorio"
	_ "consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared/apperror"
)

type ReservatorioHandler struct {
	service *reservatorioapp.Service
}

func NewReservatorioHandler(service *reservatorioapp.Service) *ReservatorioHandler {
	return &ReservatorioHandler{service: service}
}

type reservatorioRequestBody struct {
	EmpresaID     int64   `json:"empresa_id"`
	Nome          string  `json:"nome"`
	Capacidade    float64 `json:"capacidade"`
	NivelInicial  float64 `json:"nivel_inicial"`
	NivelMinimo   float64 `json:"nivel_minimo"`
	CombustivelID int64   `json:"combustivel_id"`
}

// CreateReservatorio cadastra um novo reservatório.
// @Summary Cadastrar reservatório
// @Description Cria um novo reservatório no sistema.
// @Tags Reservatórios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reservatorio body reservatorioRequestBody true "Dados do reservatório"
// @Success 201 {object} reservatorio.Reservatorio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /reservatorios [post]
func (h *ReservatorioHandler) create(w http.ResponseWriter, r *http.Request) {
	var body reservatorioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	reserv, err := h.service.Create.Handle(r.Context(), reservatorioapp.CreateCommand{
		EmpresaID:     body.EmpresaID,
		Nome:          body.Nome,
		Capacidade:    body.Capacidade,
		NivelInicial:  body.NivelInicial,
		NivelMinimo:   body.NivelMinimo,
		CombustivelID: body.CombustivelID,
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, reserv)
}

// UpdateReservatorio atualiza os dados de um reservatório existente.
// @Summary Atualizar reservatório
// @Description Atualiza os dados de um reservatório.
// @Tags Reservatórios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do reservatório"
// @Param reservatorio body reservatorioRequestBody true "Dados atualizados do reservatório"
// @Success 200 {object} reservatorio.Reservatorio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /reservatorios/{id} [put]
func (h *ReservatorioHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body reservatorioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	reserv, err := h.service.Update.Handle(r.Context(), reservatorioapp.UpdateCommand{
		ID:            id,
		Nome:          body.Nome,
		Capacidade:    body.Capacidade,
		NivelMinimo:   body.NivelMinimo,
		CombustivelID: body.CombustivelID,
		UsuarioID:     currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, reserv)
}

// DeleteReservatorio remove um reservatório do sistema.
// @Summary Excluir reservatório
// @Description Remove um reservatório do sistema.
// @Tags Reservatórios
// @Security BearerAuth
// @Param id path int true "ID do reservatório"
// @Success 204 "Reservatório excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /reservatorios/{id} [delete]
func (h *ReservatorioHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), reservatorioapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetReservatorio retorna os dados de um reservatório.
// @Summary Buscar reservatório por ID
// @Description Retorna os dados completos de um reservatório.
// @Tags Reservatórios
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do reservatório"
// @Success 200 {object} reservatorio.Reservatorio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /reservatorios/{id} [get]
func (h *ReservatorioHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	reserv, err := h.service.Get.Handle(r.Context(), reservatorioapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, reserv)
}

// ListReservatorios lista reservatórios com filtros opcionais.
// @Summary Listar reservatórios
// @Description Lista os reservatórios, podendo aplicar filtros.
// @Tags Reservatórios
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} reservatorio.Reservatorio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /reservatorios [get]
func (h *ReservatorioHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), reservatorioapp.ListQuery{
		EmpresaID:     parseQueryInt64(r, "empresa_id"),
		CombustivelID: parseQueryInt64(r, "combustivel_id"),
		Ativo:         ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
