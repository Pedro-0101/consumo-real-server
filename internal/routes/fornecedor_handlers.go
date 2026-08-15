package routes

import (
	"net/http"
	"strconv"

	fornecedorapp "consumo-real-server/internal/application/fornecedor"
	_ "consumo-real-server/internal/domain/fornecedor"
	"consumo-real-server/internal/shared/apperror"
)

type FornecedorHandler struct {
	service *fornecedorapp.Service
}

func NewFornecedorHandler(service *fornecedorapp.Service) *FornecedorHandler {
	return &FornecedorHandler{service: service}
}

type fornecedorRequestBody struct {
	EmpresaID int64  `json:"empresa_id"`
	Nome      string `json:"nome"`
	CNPJ      string `json:"cnpj"`
}

// CreateFornecedor cadastra um novo fornecedor.
// @Summary Cadastrar fornecedor
// @Description Cria um novo fornecedor no sistema.
// @Tags Fornecedores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fornecedor body fornecedorRequestBody true "Dados do fornecedor"
// @Success 201 {object} fornecedor.Fornecedor
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /fornecedores [post]
func (h *FornecedorHandler) create(w http.ResponseWriter, r *http.Request) {
	var body fornecedorRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Create.Handle(r.Context(), fornecedorapp.CreateCommand{
		EmpresaID: body.EmpresaID,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, f)
}

// UpdateFornecedor atualiza os dados de um fornecedor existente.
// @Summary Atualizar fornecedor
// @Description Atualiza os dados de um fornecedor.
// @Tags Fornecedores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do fornecedor"
// @Param fornecedor body fornecedorRequestBody true "Dados atualizados do fornecedor"
// @Success 200 {object} fornecedor.Fornecedor
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /fornecedores/{id} [put]
func (h *FornecedorHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body fornecedorRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	f, err := h.service.Update.Handle(r.Context(), fornecedorapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

// DeleteFornecedor remove um fornecedor do sistema.
// @Summary Excluir fornecedor
// @Description Remove um fornecedor do sistema.
// @Tags Fornecedores
// @Security BearerAuth
// @Param id path int true "ID do fornecedor"
// @Success 204 "Fornecedor excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /fornecedores/{id} [delete]
func (h *FornecedorHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), fornecedorapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetFornecedor retorna os dados de um fornecedor.
// @Summary Buscar fornecedor por ID
// @Description Retorna os dados completos de um fornecedor.
// @Tags Fornecedores
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do fornecedor"
// @Success 200 {object} fornecedor.Fornecedor
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /fornecedores/{id} [get]
func (h *FornecedorHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	f, err := h.service.Get.Handle(r.Context(), fornecedorapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, f)
}

// ListFornecedores lista fornecedores com filtros opcionais.
// @Summary Listar fornecedores
// @Description Lista os fornecedores, podendo aplicar filtros.
// @Tags Fornecedores
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} fornecedor.Fornecedor
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /fornecedores [get]
func (h *FornecedorHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), fornecedorapp.ListQuery{
		EmpresaID: parseQueryInt64(r, "empresa_id"),
		CNPJ:      r.URL.Query().Get("cnpj"),
		Ativo:     ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
