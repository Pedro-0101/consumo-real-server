package frentista

import (
	domainfrentista "consumo-real-server/internal/domain/frentista"
)

// Service agrupa os handlers de frentista.
type Service struct {
	Create          *CreateHandler
	Update          *UpdateHandler
	Delete          *DeleteHandler
	Get             *GetHandler
	List            *ListHandler
	VincularUsuario *VincularUsuarioHandler
}

func NewService(repo domainfrentista.Repository) *Service {
	return &Service{
		Create:          NewCreateHandler(repo),
		Update:          NewUpdateHandler(repo),
		Delete:          NewDeleteHandler(repo),
		Get:             NewGetHandler(repo),
		List:            NewListHandler(repo),
		VincularUsuario: NewVincularUsuarioHandler(repo),
	}
}
