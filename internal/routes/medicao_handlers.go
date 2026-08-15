package routes

import (
	"net/http"

	medicaoapp "consumo-real-server/internal/application/medicao"
	_ "consumo-real-server/internal/domain/medicao"
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

// CreateMedicao registra uma nova medição.
// @Summary Cadastrar Medição
// @Description Registra uma nova medição de reservatório.
// @Tags Medições
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param medicao body medicaoRequestBody true "Dados da medição"
// @Success 201 {object} medicao.Medicao
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /medicoes [post]
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

// UpdateMedicao atualiza os dados de uma medição existente.
// @Summary Atualizar Medição
// @Description Atualiza os dados de uma medição.
// @Tags Medições
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da medição"
// @Param medicao body medicaoRequestBody true "Dados atualizados da medição"
// @Success 200 {object} medicao.Medicao
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /medicoes/{id} [put]
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

// DeleteMedicao remove uma medição do sistema.
// @Summary Excluir Medição
// @Description Remove uma medição do sistema.
// @Tags Medições
// @Security BearerAuth
// @Param id path int true "ID da medição"
// @Success 204 "Medição excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /medicoes/{id} [delete]
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

// GetMedicao retorna os dados de uma medição.
// @Summary Buscar Medição por ID
// @Description Retorna os dados completos de uma medição.
// @Tags Medições
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da medição"
// @Success 200 {object} medicao.Medicao
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /medicoes/{id} [get]
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

// ListMedicoes lista medições com filtros opcionais.
// @Summary Listar Medições
// @Description Lista as medições, podendo aplicar filtros.
// @Tags Medições
// @Produce json
// @Security BearerAuth
// @Param empresa_id query int false "ID da empresa"
// @Param reservatorio_id query int false "ID do reservatório"
// @Success 200 {array} medicao.Medicao
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /medicoes [get]
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
