package fornecedor

import (
	domainfornecedor "consumo-real-server/internal/domain/fornecedor"
)

// Service agrupa os handlers de fornecedor.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainfornecedor.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
