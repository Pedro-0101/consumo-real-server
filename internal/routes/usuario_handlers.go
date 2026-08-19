package routes

import (
	"net/http"
	"strconv"

	usuarioapp "consumo-real-server/internal/application/usuario"
	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared/apperror"
)

type UsuarioHandler struct {
	service *usuarioapp.Service
}

func NewUsuarioHandler(service *usuarioapp.Service) *UsuarioHandler {
	return &UsuarioHandler{service: service}
}

type usuarioRequestBody struct {
	EmpresaID int64  `json:"empresa_id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Senha     string `json:"senha"`
	Papel     string `json:"papel"`
}

// CreateUsuario cadastra um novo usuário.
// @Summary Cadastrar usuário
// @Description Cria um novo usuário no sistema.
// @Tags Usuários
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param usuario body usuarioRequestBody true "Dados do usuário"
// @Success 201 {object} domainusuario.Usuario
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios [post]
func (h *UsuarioHandler) create(w http.ResponseWriter, r *http.Request) {
	var body usuarioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	u, err := h.service.Create.Handle(r.Context(), usuarioapp.CreateCommand{
		EmpresaID: body.EmpresaID,
		Nome:      body.Nome,
		Email:     body.Email,
		Senha:     body.Senha,
		Papel:     domainusuario.Papel(body.Papel),
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, u)
}

// UpdateUsuario atualiza os dados de um usuário existente.
// @Summary Atualizar usuário
// @Description Atualiza nome, e-mail e papel de um usuário.
// @Tags Usuários
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Param usuario body usuarioRequestBody true "Dados atualizados do usuário"
// @Success 200 {object} domainusuario.Usuario
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios/{id} [put]
func (h *UsuarioHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body usuarioRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	u, err := h.service.Update.Handle(r.Context(), usuarioapp.UpdateCommand{
		ID:        id,
		EmpresaID: body.EmpresaID,
		Nome:      body.Nome,
		Email:     body.Email,
		Papel:     domainusuario.Papel(body.Papel),
		UsuarioID: currentUserID(r),
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}

type changePasswordRequestBody struct {
	SenhaAtual string `json:"senha_atual"`
	NovaSenha  string `json:"nova_senha"`
}

// ChangePassword altera a senha do usuário informando a senha atual.
// @Summary Alterar senha
// @Description Altera a senha do usuário autenticado.
// @Tags Usuários
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Param senha body changePasswordRequestBody true "Senha atual e nova senha"
// @Success 204 "Senha alterada com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios/{id}/senha [patch]
func (h *UsuarioHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body changePasswordRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	if err := h.service.ChangePassword.Handle(r.Context(), usuarioapp.ChangePasswordCommand{
		ID:         id,
		SenhaAtual: body.SenhaAtual,
		NovaSenha:  body.NovaSenha,
		UsuarioID:  currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUsuario remove um usuário do sistema.
// @Summary Excluir usuário
// @Description Remove logicamente um usuário do sistema.
// @Tags Usuários
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 204 "Usuário excluído com sucesso"
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios/{id} [delete]
func (h *UsuarioHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	if err := h.service.Delete.Handle(r.Context(), usuarioapp.DeleteCommand{
		ID:        id,
		UsuarioID: currentUserID(r),
	}); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetUsuario retorna os dados de um usuário.
// @Summary Buscar usuário por ID
// @Description Retorna os dados completos de um usuário.
// @Tags Usuários
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} domainusuario.Usuario
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 404 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios/{id} [get]
func (h *UsuarioHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	u, err := h.service.Get.Handle(r.Context(), usuarioapp.GetQuery{ID: id})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}

// ListUsuarios lista usuários com filtros opcionais.
// @Summary Listar usuários
// @Description Lista os usuários, podendo filtrar por papel e status ativo.
// @Tags Usuários
// @Produce json
// @Security BearerAuth
// @Param papel query string false "Papel do usuário (ex.: GERENTE)"
// @Param ativo query bool false "Filtrar apenas usuários ativos"
// @Success 200 {array} domainusuario.Usuario
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /usuarios [get]
func (h *UsuarioHandler) list(w http.ResponseWriter, r *http.Request) {
	var ativo *bool
	if raw := r.URL.Query().Get("ativo"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperror.WriteError(w, apperror.Validation("parâmetro 'ativo' inválido", err))
			return
		}
		ativo = &parsed
	}

	empresaID := currentUserEmpresaID(r)
	if empresaID <= 0 {
		empresaID = parseQueryInt64(r, "empresa_id")
	}

	list, err := h.service.List.Handle(r.Context(), usuarioapp.ListQuery{
		EmpresaID: empresaID,
		Papel:     domainusuario.Papel(r.URL.Query().Get("papel")),
		Ativo:     ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

// Me retorna o usuário autenticado pelo token.
// @Summary Usuário autenticado
// @Description Retorna os dados do usuário autenticado com base no token JWT.
// @Tags Usuários
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domainusuario.Usuario
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /auth/me [get]
func (h *UsuarioHandler) me(w http.ResponseWriter, r *http.Request) {
	u, err := h.service.Me.Handle(r.Context(), usuarioapp.MeQuery{UsuarioID: currentUserID(r)})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}
