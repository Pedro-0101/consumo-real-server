package ordemabastecimento

import (
	domainordem "consumo-real-server/internal/domain/ordemabastecimento"
)

// Service agrupa os handlers de ordem de abastecimento.
type Service struct {
	Create    *CreateHandler
	Update    *UpdateHandler
	Delete    *DeleteHandler
	Get       *GetHandler
	List      *ListHandler
	Autorizar *AutorizarHandler
}

func NewService(repo domainordem.Repository) *Service {
	return &Service{
		Create:    NewCreateHandler(repo),
		Update:    NewUpdateHandler(repo),
		Delete:    NewDeleteHandler(repo),
		Get:       NewGetHandler(repo),
		List:      NewListHandler(repo),
		Autorizar: NewAutorizarHandler(repo),
	}
}
