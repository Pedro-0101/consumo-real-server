package medicao

import (
	domainmedicao "consumo-real-server/internal/domain/medicao"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
)

// Service agrupa os handlers de medição.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainmedicao.Repository, reservRepo domainreservatorio.Repository) *Service {
	return &Service{
		Create: NewCreateHandler(repo, reservRepo),
		Update: NewUpdateHandler(repo, reservRepo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
