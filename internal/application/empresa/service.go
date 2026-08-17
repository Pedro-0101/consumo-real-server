package empresa

import (
	domainempresa "consumo-real-server/internal/domain/empresa"
	domainusuario "consumo-real-server/internal/domain/usuario"
)

// Service agrupa os handlers de empresa e é o ponto de entrada
// da camada de aplicação para a camada de apresentação.
type Service struct {
	Create *CreateHandler
	Update *UpdateHandler
	Delete *DeleteHandler
	Get    *GetHandler
	List   *ListHandler
}

func NewService(repo domainempresa.Repository, onboarding OnboardingRepository, hasher domainusuario.PasswordHasher) *Service {
	return &Service{
		Create: NewCreateHandler(repo, onboarding, hasher),
		Update: NewUpdateHandler(repo),
		Delete: NewDeleteHandler(repo),
		Get:    NewGetHandler(repo),
		List:   NewListHandler(repo),
	}
}
