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

// CreatePatrimonio cadastra um novo patrimônio.
// @Summary Cadastrar Patrimônio
// @Description Cria um novo patrimônio no sistema.
// @Tags Patrimônios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patrimonio body patrimonioRequestBody true "Dados do patrimônio"
// @Success 201 {object} domainpatrimonio.Patrimonio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /patrimonios [post]
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

// UpdatePatrimonio atualiza os dados de um patrimônio existente.
// @Summary Atualizar Patrimônio
// @Description Atualiza os dados de um patrimônio.
// @Tags Patrimônios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do patrimônio"
// @Param patrimonio body patrimonioRequestBody true "Dados atualizados do patrimônio"
// @Success 200 {object} domainpatrimonio.Patrimonio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /patrimonios/{id} [put]
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

// DeletePatrimonio remove um patrimônio do sistema.
// @Summary Excluir Patrimônio
// @Description Remove um patrimônio do sistema.
// @Tags Patrimônios
// @Security BearerAuth
// @Param id path int true "ID do patrimônio"
// @Success 204 "Patrimônio excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /patrimonios/{id} [delete]
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

// GetPatrimonio retorna os dados de um patrimônio.
// @Summary Buscar Patrimônio por ID
// @Description Retorna os dados completos de um patrimônio.
// @Tags Patrimônios
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do patrimônio"
// @Success 200 {object} domainpatrimonio.Patrimonio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /patrimonios/{id} [get]
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

// ListPatrimonios lista patrimônios com filtros opcionais.
// @Summary Listar Patrimônios
// @Description Lista os patrimônios, podendo aplicar filtros.
// @Tags Patrimônios
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} domainpatrimonio.Patrimonio
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /patrimonios [get]
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
