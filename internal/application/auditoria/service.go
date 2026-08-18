package auditoria

import (
	domainauditoria "consumo-real-server/internal/domain/auditoria"
)

// Service agrupa os handlers de auditoria e é o ponto de entrada
// da camada de aplicação para a camada de apresentação.
type Service struct {
	List *ListHandler
}

func NewService(repo domainauditoria.Repository) *Service {
	return &Service{
		List: NewListHandler(repo),
	}
}