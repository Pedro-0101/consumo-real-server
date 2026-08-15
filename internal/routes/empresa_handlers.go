package routes

import (
	"net/http"
	"strconv"

	empresaapp "consumo-real-server/internal/application/empresa"
	_ "consumo-real-server/internal/domain/empresa"
	"consumo-real-server/internal/shared/apperror"
)

type EmpresaHandler struct {
	service *empresaapp.Service
}

func NewEmpresaHandler(service *empresaapp.Service) *EmpresaHandler {
	return &EmpresaHandler{service: service}
}

type empresaRequestBody struct {
	Nome string `json:"nome"`
	CNPJ string `json:"cnpj"`
}

// CreateEmpresa cadastra uma nova empresa.
// @Summary Cadastrar Empresa
// @Description Cria uma nova empresa no sistema.
// @Tags Empresas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param empresa body empresaRequestBody true "Dados da empresa"
// @Success 201 {object} empresa.Empresa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /empresas [post]
func (h *EmpresaHandler) create(w http.ResponseWriter, r *http.Request) {
	var body empresaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Create.Handle(r.Context(), empresaapp.CreateCommand{
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, e)
}

// UpdateEmpresa atualiza os dados de uma empresa existente.
// @Summary Atualizar Empresa
// @Description Atualiza os dados de uma empresa.
// @Tags Empresas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da empresa"
// @Param empresa body empresaRequestBody true "Dados atualizados da empresa"
// @Success 200 {object} empresa.Empresa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /empresas/{id} [put]
func (h *EmpresaHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body empresaRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	e, err := h.service.Update.Handle(r.Context(), empresaapp.UpdateCommand{
		ID:        id,
		Nome:      body.Nome,
		CNPJ:      body.CNPJ,
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

// DeleteEmpresa remove uma empresa do sistema.
// @Summary Excluir Empresa
// @Description Remove uma empresa do sistema.
// @Tags Empresas
// @Security BearerAuth
// @Param id path int true "ID da empresa"
// @Success 204 "Empresa excluída com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /empresas/{id} [delete]
func (h *EmpresaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), empresaapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetEmpresa retorna os dados de uma empresa.
// @Summary Buscar Empresa por ID
// @Description Retorna os dados completos de uma empresa.
// @Tags Empresas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da empresa"
// @Success 200 {object} empresa.Empresa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /empresas/{id} [get]
func (h *EmpresaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	e, err := h.service.Get.Handle(r.Context(), empresaapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, e)
}

// ListEmpresas lista empresas com filtros opcionais.
// @Summary Listar Empresas
// @Description Lista as empresas, podendo aplicar filtros.
// @Tags Empresas
// @Produce json
// @Security BearerAuth
// @Param ativo query bool false "Filtrar apenas registros ativos"
// @Success 200 {array} empresa.Empresa
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /empresas [get]
func (h *EmpresaHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	list, err := h.service.List.Handle(r.Context(), empresaapp.ListQuery{Ativo: ativo})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}
