package patrimonio

import (
	domainpatrimonio "consumo-real-server/internal/domain/patrimonio"
)

// Service agrupa os handlers de patrimônio.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainpatrimonio.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
