package routes

import (
	"net/http"
	"strconv"
	"time"

	precoapp "consumo-real-server/internal/application/preco"
	_ "consumo-real-server/internal/domain/preco"
	"consumo-real-server/internal/shared/apperror"
)

type PrecoHandler struct {
	service *precoapp.Service
}

func NewPrecoHandler(service *precoapp.Service) *PrecoHandler {
	return &PrecoHandler{service: service}
}

type precoRequestBody struct {
	EmpresaID      int64   `json:"empresa_id"`
	CombustivelID  int64   `json:"combustivel_id"`
	PrecoCusto     float64 `json:"preco_custo"`
	PrecoVenda     float64 `json:"preco_venda"`
	VigenciaInicio string  `json:"vigencia_inicio"`
}

func (h *PrecoHandler) parseVigenciaInicio(w http.ResponseWriter, raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, true
	}
	vigencia, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		apperror.WriteError(w, apperror.Validation("parâmetro 'vigencia_inicio' inválido, use RFC3339 (ex.: 2026-01-01T00:00:00Z)", err))
		return time.Time{}, false
	}
	return vigencia, true
}

// CreatePreco cadastra um novo preço.
// @Summary Cadastrar preço
// @Description Cria um novo preço no sistema.
// @Tags Preços
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param preco body precoRequestBody true "Dados do preço"
// @Success 201 {object} preco.Preco
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /precos [post]
func (h *PrecoHandler) create(w http.ResponseWriter, r *http.Request) {
	var body precoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	vigencia, ok := h.parseVigenciaInicio(w, body.VigenciaInicio)
	if !ok {
		return
	}

	p, err := h.service.Create.Handle(r.Context(), precoapp.CreateCommand{
		EmpresaID:      body.EmpresaID,
		CombustivelID:  body.CombustivelID,
		PrecoCusto:     body.PrecoCusto,
		PrecoVenda:     body.PrecoVenda,
		VigenciaInicio: vigencia,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, p)
}

// UpdatePreco atualiza os dados de um preço existente.
// @Summary Atualizar preço
// @Description Atualiza os dados de um preço.
// @Tags Preços
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do preço"
// @Param preco body precoRequestBody true "Dados atualizados do preço"
// @Success 200 {object} preco.Preco
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /precos/{id} [put]
func (h *PrecoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body precoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	vigencia, ok := h.parseVigenciaInicio(w, body.VigenciaInicio)
	if !ok {
		return
	}

	p, err := h.service.Update.Handle(r.Context(), precoapp.UpdateCommand{
		ID:             id,
		PrecoCusto:     body.PrecoCusto,
		PrecoVenda:     body.PrecoVenda,
		VigenciaInicio: vigencia,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, p)
}

// DeletePreco remove um preço do sistema.
// @Summary Excluir preço
// @Description Remove um preço do sistema.
// @Tags Preços
// @Security BearerAuth
// @Param id path int true "ID do preço"
// @Success 204 "Preço excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /precos/{id} [delete]
func (h *PrecoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), precoapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetPreco retorna os dados de um preço.
// @Summary Buscar preço por ID
// @Description Retorna os dados completos de um preço.
// @Tags Preços
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do preço"
// @Success 200 {object} preco.Preco
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /precos/{id} [get]
func (h *PrecoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	p, err := h.service.Get.Handle(r.Context(), precoapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, p)
}

// ListPrecos lista preços com filtros opcionais.
// @Summary Listar preços
// @Description Lista os preços, podendo aplicar filtros.
// @Tags Preços
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} preco.Preco
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /precos [get]
func (h *PrecoHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), precoapp.ListQuery{
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
