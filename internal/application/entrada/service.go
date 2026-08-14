package entrada

import (
	domainentrada "consumo-real-server/internal/domain/entrada"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
)

// Service agrupa os handlers de entrada.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainentrada.Repository, reservRepo domainreservatorio.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo, reservRepo),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
