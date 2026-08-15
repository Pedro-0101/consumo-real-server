package routes

import (
	"net/http"

	usuarioapp "consumo-real-server/internal/application/usuario"
	"consumo-real-server/internal/shared/apperror"
)

type AuthHandler struct {
	service *usuarioapp.Service
}

func NewAuthHandler(service *usuarioapp.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequestBody struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Login efetua a autenticação do usuário e retorna o token JWT.
// @Summary Autenticar usuário
// @Description Autentica o usuário com e-mail e senha, retornando um token JWT de acesso.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param login body loginRequestBody true "Credenciais de acesso"
// @Success 200 {object} usuarioapp.LoginResult
// @Failure 400 {object} apperror.ErrorResponse
// @Failure 401 {object} apperror.ErrorResponse
// @Failure 422 {object} apperror.ErrorResponse
// @Failure 500 {object} apperror.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody
	if err := apperror.DecodeJSON(w, r, &body); err != nil {
		apperror.WriteError(w, err)
		return
	}

	result, err := h.service.Login.Handle(r.Context(), usuarioapp.LoginCommand{
		Email: body.Email,
		Senha: body.Senha,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, result)
}
