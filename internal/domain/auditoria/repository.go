package auditoria

import (
	"context"
	"time"
)

// ListFilter define os filtros disponíveis para a consulta de auditorias.
type ListFilter struct {
	EmpresaID int64
	Entidade  string
	Operacao  Operacao
	UsuarioID int64
	De        *time.Time
	Ate       *time.Time
	Limit     int
	Offset    int
}

// Repository é o contrato de persistência do agregado Auditoria.
type Repository interface {
	Create(ctx context.Context, a *Auditoria) error
	List(ctx context.Context, filter ListFilter) ([]Auditoria, error)
}