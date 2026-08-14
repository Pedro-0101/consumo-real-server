package apperror

import (
	"errors"
	"fmt"
)

// Kind classifica o tipo de erro padronizado do sistema.
type Kind string

const (
	KindValidation   Kind = "VALIDATION"
	KindNotFound     Kind = "NOT_FOUND"
	KindConflict     Kind = "CONFLICT"
	KindUnauthorized Kind = "UNAUTHORIZED"
	KindInternal     Kind = "INTERNAL"
)

// AppError é o erro padronizado da aplicação. Erros de domínio são
// convertidos para AppError na camada de aplicação.
type AppError struct {
	Kind    Kind
	Message string
	Details any
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Cause }

// Validation cria um erro de validação (regra de negócio/entrada inválida).
func Validation(message string, cause error) *AppError {
	return &AppError{Kind: KindValidation, Message: message, Cause: cause}
}

// NotFound cria um erro de recurso não encontrado.
func NotFound(message string) *AppError {
	return &AppError{Kind: KindNotFound, Message: message}
}

// Conflict cria um erro de conflito de estado (ex.: registro duplicado).
func Conflict(message string, cause error) *AppError {
	return &AppError{Kind: KindConflict, Message: message, Cause: cause}
}

// Unauthorized cria um erro de autenticação/autorização.
func Unauthorized(message string) *AppError {
	return &AppError{Kind: KindUnauthorized, Message: message}
}

// Internal cria um erro interno não esperado.
func Internal(message string, cause error) *AppError {
	return &AppError{Kind: KindInternal, Message: message, Cause: cause}
}

// As extrai o AppError da cadeia de erros, se existir.
func As(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// FromDomain converte erros de domínio (regras de negócio) em erro padronizado.
// Erros que já são AppError são retornados como estão.
func FromDomain(err error) *AppError {
	if ae, ok := As(err); ok {
		return ae
	}
	return &AppError{Kind: KindValidation, Message: err.Error(), Cause: err}
}
