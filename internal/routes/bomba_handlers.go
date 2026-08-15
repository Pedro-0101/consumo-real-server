package routes

import (
	"net/http"
	"strconv"

	bombaapp "consumo-real-server/internal/application/bomba"
	_ "consumo-real-server/internal/domain/bomba"
	"consumo-real-server/internal/shared/apperror"
)

type BombaHandler struct {
	service *bombaapp.Service
}

func NewBombaHandler(service *bombaapp.Service) *BombaHandler {
	return &BombaHandler{service: service}
}

type bombaRequestBody struct {
	EmpresaID      int64  `json:"empresa_id"`
	LocalID        int64  `json:"local_id"`
	ReservatorioID int64  `json:"reservatorio_id"`
	Movel          bool   `json:"movel"`
	Nome           string `json:"nome"`
	Descricao      string `json:"descricao"`
}

// CreateBomba cadastra uma nova bomba.
// @Summary Cadastrar bomba
// @Description Cria uma nova bomba no sistema.
// @Tags Bombas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bomba body bombaRequestBody true "Dados da bomba"
// @Success 201 {object} bomba.Bomba
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas [post]
func (h *BombaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body bombaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	bomba, err := h.service.Create.Handle(r.Context(), bombaapp.CreateCommand{
		EmpresaID:      body.EmpresaID,
		LocalID:        body.LocalID,
		ReservatorioID: body.ReservatorioID,
		Movel:          body.Movel,
		Nome:           body.Nome,
		Descricao:      body.Descricao,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, bomba)
}

// UpdateBomba atualiza os dados de uma bomba existente.
// @Summary Atualizar bomba
// @Description Atualiza os dados de uma bomba.
// @Tags Bombas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da bomba"
// @Param bomba body bombaRequestBody true "Dados atualizados da bomba"
// @Success 200 {object} bomba.Bomba
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas/{id} [put]
func (h *BombaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body bombaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	bomba, err := h.service.Update.Handle(r.Context(), bombaapp.UpdateCommand{
		ID:             id,
		LocalID:        body.LocalID,
		ReservatorioID: body.ReservatorioID,
		Movel:          body.Movel,
		Nome:           body.Nome,
		Descricao:      body.Descricao,
		UsuarioID:      currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, bomba)
}

// DeleteBomba remove uma bomba do sistema.
// @Summary Excluir bomba
// @Description Remove uma bomba do sistema.
// @Tags Bombas
// @Security BearerAuth
// @Param id path int true "ID da bomba"
// @Success 204 "Bomba excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas/{id} [delete]
func (h *BombaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), bombaapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetBomba retorna os dados de uma bomba.
// @Summary Buscar bomba por ID
// @Description Retorna os dados completos de uma bomba.
// @Tags Bombas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da bomba"
// @Success 200 {object} bomba.Bomba
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas/{id} [get]
func (h *BombaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	bomba, err := h.service.Get.Handle(r.Context(), bombaapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, bomba)
}

// ListBombas lista bombas com filtros opcionais.
// @Summary Listar bombas
// @Description Lista as bombas, podendo aplicar filtros.
// @Tags Bombas
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} bomba.Bomba
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas [get]
func (h *BombaHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), bombaapp.ListQuery{
		EmpresaID:      parseQueryInt64(r, "empresa_id"),
		LocalID:        parseQueryInt64(r, "local_id"),
		ReservatorioID: parseQueryInt64(r, "reservatorio_id"),
		Ativo:          ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

type bicoRequestBody struct {
	Nome string `json:"nome"`
}

// AdicionarBico adiciona um bico a uma bomba.
// @Summary Adicionar bico
// @Description Adiciona um novo bico a uma bomba existente.
// @Tags Bombas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da bomba"
// @Param bico body bicoRequestBody true "Dados do bico"
// @Success 201 {object} bomba.Bico
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas/{id}/bicos [post]
func (h *BombaHandler) adicionarBico(w http.ResponseWriter, r *http.Request) {
	bombaID, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body bicoRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	bico, err := h.service.AdicionarBico.Handle(r.Context(), bombaapp.AdicionarBicoCommand{
		BombaID:   bombaID,
		Nome:      body.Nome,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, bico)
}

// DesativarBico desativa um bico de uma bomba.
// @Summary Desativar bico
// @Description Desativa um bico de uma bomba existente.
// @Tags Bombas
// @Security BearerAuth
// @Param id path int true "ID da bomba"
// @Param bico_id query int true "ID do bico"
// @Success 204 "Bico desativado com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /bombas/{id}/bicos [delete]
func (h *BombaHandler) desativarBico(w http.ResponseWriter, r *http.Request) {
	bombaID, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	bicoID, err := strconv.ParseInt(r.URL.Query().Get("bico_id"), 10, 64)
	if err != nil || bicoID <= 0 {
		apperror.WriteError(w, apperror.Validation("parâmetro 'bico_id' inválido", err))
		return
	}

	if err := h.service.DesativarBico.Handle(r.Context(), bombaapp.DesativarBicoCommand{
		BombaID:   bombaID,
		BicoID:    bicoID,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
