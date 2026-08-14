package bomba

import (
	domainbomba "consumo-real-server/internal/domain/bomba"
)

// Service agrupa os handlers de bomba e bicos.
type Service struct {
	Create        *CreateHandler
	Update        *UpdateHandler
	Delete        *DeleteHandler
	Get           *GetHandler
	List          *ListHandler
	AdicionarBico *AdicionarBicoHandler
	DesativarBico *DesativarBicoHandler
}

func NewService(repo domainbomba.Repository) *Service {
	return &Service{
		Create:        NewCreateHandler(repo),
		Update:        NewUpdateHandler(repo),
		Delete:        NewDeleteHandler(repo),
		Get:           NewGetHandler(repo),
		List:          NewListHandler(repo),
		AdicionarBico: NewAdicionarBicoHandler(repo),
		DesativarBico: NewDesativarBicoHandler(repo),
	}
}
