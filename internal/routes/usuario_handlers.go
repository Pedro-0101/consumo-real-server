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

func (h *UsuarioHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		apperror.WriteError(w, apperror.Validation("id inválido", nil))
		return
	}

	var body struct {
		SenhaAtual string `json:"senha_atual"`
		NovaSenha  string `json:"nova_senha"`
	}
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

	list, err := h.service.List.Handle(r.Context(), usuarioapp.ListQuery{
		Papel: domainusuario.Papel(r.URL.Query().Get("papel")),
		Ativo: ativo,
	})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, list)
}

func (h *UsuarioHandler) me(w http.ResponseWriter, r *http.Request) {
	u, err := h.service.Me.Handle(r.Context(), usuarioapp.MeQuery{UsuarioID: currentUserID(r)})
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, u)
}
