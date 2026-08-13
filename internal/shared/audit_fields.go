package shared

import (
	"time"

	"consumo-real-server/internal/shared/timeutil"
)

/*

Os campos de auditoria são utilizados para registrar informações sobre a criação e atualização de registros no sistema.
Eles incluem a data e hora de criação e atualização, bem como o ID do usuário responsável por essas ações.

*/

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int64
	UpdatedBy int64
}

func NewAuditFields(createdBy int64) AuditFields {
	return AuditFields{
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
}

func (a *AuditFields) Update(updatedBy int64) {
	a.UpdatedAt = timeutil.Now()
	a.UpdatedBy = updatedBy
}
