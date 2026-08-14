package reservatorio

import (
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
)

// Service agrupa os handlers de reservatório.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainreservatorio.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
