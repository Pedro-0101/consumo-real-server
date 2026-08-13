package seeds

import (
	"errors"
	"strings"

	"consumo-real-server/internal/domain/usuario"
)

type AdminBaseRepository interface {
	ExistsByEmail(email string) (bool, error)
	Save(u *usuario.Usuario) error
}

var ErrAdminBaseInvalido = errors.New("dados do usuário admin base inválidos")

// SeedAdminBase garante que o usuário administrador base do sistema sempre exista.
// Deve ser executado na inicialização da aplicação.
func SeedAdminBase(repo AdminBaseRepository, nome, email, senhaHash string) error {
	exists, err := repo.ExistsByEmail(strings.ToLower(email))
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	admin, err := usuario.NewAdminBase(nome, email, senhaHash)
	if err != nil {
		return ErrAdminBaseInvalido
	}
	return repo.Save(admin)
}
