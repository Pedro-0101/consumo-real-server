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

// CreateCombustivel cadastra um novo combustível.
// @Summary Cadastrar Combustível
// @Description Cria um novo combustível no sistema.
// @Tags Combustíveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param combustivel body combustivelRequestBody true "Dados do combustível"
// @Success 201 {object} domaincombustivel.Combustivel
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /combustiveis [post]
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

// UpdateCombustivel atualiza os dados de um combustível existente.
// @Summary Atualizar Combustível
// @Description Atualiza os dados de um combustível.
// @Tags Combustíveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do combustível"
// @Param combustivel body combustivelRequestBody true "Dados atualizados do combustível"
// @Success 200 {object} domaincombustivel.Combustivel
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /combustiveis/{id} [put]
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

// DeleteCombustivel remove um combustível do sistema.
// @Summary Excluir Combustível
// @Description Remove um combustível do sistema.
// @Tags Combustíveis
// @Security BearerAuth
// @Param id path int true "ID do combustível"
// @Success 204 "Combustível excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /combustiveis/{id} [delete]
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

// GetCombustivel retorna os dados de um combustível.
// @Summary Buscar Combustível por ID
// @Description Retorna os dados completos de um combustível.
// @Tags Combustíveis
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do combustível"
// @Success 200 {object} domaincombustivel.Combustivel
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /combustiveis/{id} [get]
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

// ListCombustiveis lista combustíveis com filtros opcionais.
// @Summary Listar Combustíveis
// @Description Lista os combustíveis, podendo aplicar filtros.
// @Tags Combustíveis
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} domaincombustivel.Combustivel
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /combustiveis [get]
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
