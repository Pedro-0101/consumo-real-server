package preco

import (
	domainpreco "consumo-real-server/internal/domain/preco"
)

// Service agrupa os handlers de preço.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainpreco.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
