package abastecimento

import (
	domainabastecimento "consumo-real-server/internal/domain/abastecimento"
	domainbomba "consumo-real-server/internal/domain/bomba"
	domainreservatorio "consumo-real-server/internal/domain/reservatorio"
)

// Service agrupa os handlers de abastecimento.
type Service struct {
	Create         *CreateHandler
	CreateTransfer *CreateTransferenciaHandler
	Update         *UpdateHandler
	Delete         *DeleteHandler
	Get            *GetHandler
	List           *ListHandler
}

func NewService(repo domainabastecimento.Repository, bombaRepo domainbomba.Repository, reservRepo domainreservatorio.Repository) *Service {
	return &Service{
		Create:         NewCreateHandler(repo, bombaRepo, reservRepo),
		CreateTransfer: NewCreateTransferenciaHandler(repo, reservRepo),
		Update:         NewUpdateHandler(repo),
		Delete:         NewDeleteHandler(repo),
		Get:            NewGetHandler(repo),
		List:           NewListHandler(repo),
	}
}
