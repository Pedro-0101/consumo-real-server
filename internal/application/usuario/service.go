package usuario

import (
	domainusuario "consumo-real-server/internal/domain/usuario"
	"consumo-real-server/internal/shared/auth"
)

// Service agrupa os handlers de usuário e é o ponto de entrada
// da camada de aplicação para a camada de apresentação.
type Service struct {
	Create         *CreateHandler
	Update         *UpdateHandler
	ChangePassword *ChangePasswordHandler
	Delete         *DeleteHandler
	Get            *GetHandler
	List           *ListHandler
	Login          *LoginHandler
	Me             *MeHandler
}

func NewService(repo domainusuario.Repository, hasher domainusuario.PasswordHasher, tokens auth.TokenManager) *Service {
	return &Service{
		Create:         NewCreateHandler(repo, hasher),
		Update:         NewUpdateHandler(repo),
		ChangePassword: NewChangePasswordHandler(repo, hasher),
		Delete:         NewDeleteHandler(repo),
		Get:            NewGetHandler(repo),
		List:           NewListHandler(repo),
		Login:          NewLoginHandler(repo, hasher, tokens),
		Me:             NewMeHandler(repo),
	}
}
