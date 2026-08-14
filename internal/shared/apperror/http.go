package apperror

import (
	"encoding/json"
	"net/http"
)

var statusByKind = map[Kind]int{
	KindValidation:   http.StatusUnprocessableEntity,
	KindNotFound:     http.StatusNotFound,
	KindConflict:     http.StatusConflict,
	KindUnauthorized: http.StatusUnauthorized,
	KindInternal:     http.StatusInternalServerError,
}

// StatusCode retorna o código HTTP correspondente ao erro padronizado.
func StatusCode(e *AppError) int {
	if status, ok := statusByKind[e.Kind]; ok {
		return status
	}
	return http.StatusInternalServerError
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Kind    Kind   `json:"kind"`
	Message string `json:"message"`
}

// WriteError serializa qualquer erro como resposta JSON padronizada.
// Erros que não são AppError são tratados como erro interno.
func WriteError(w http.ResponseWriter, err error) {
	ae, ok := As(err)
	if !ok {
		ae = &AppError{Kind: KindInternal, Message: "erro interno", Cause: err}
	}
	WriteJSON(w, StatusCode(ae), errorResponse{
		Error: errorBody{Kind: ae.Kind, Message: ae.Message},
	})
}

// WriteJSON serializa o corpo como JSON com o status informado.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// DecodeJSON decodifica o corpo da requisição em dst, retornando erro
// padronizado de validação em caso de corpo malformado.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return Validation("corpo da requisição inválido", err)
	}
	return nil
}
