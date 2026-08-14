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

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
		Senha string `json:"senha"`
	}
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
