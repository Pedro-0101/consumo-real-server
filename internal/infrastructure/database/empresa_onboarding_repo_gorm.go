package database

import (
	"context"

	"gorm.io/gorm"

	domainempresa "consumo-real-server/internal/domain/empresa"
	domainusuario "consumo-real-server/internal/domain/usuario"
)

// EmpresaOnboardingGORMRepository persiste a empresa e o seu primeiro
// administrador em uma única transação.
type EmpresaOnboardingGORMRepository struct {
	db *gorm.DB
}

func NewEmpresaOnboardingGORMRepository(db *gorm.DB) *EmpresaOnboardingGORMRepository {
	return &EmpresaOnboardingGORMRepository{db: db}
}

func (r *EmpresaOnboardingGORMRepository) CriarEmpresaComAdministrador(
	ctx context.Context,
	e *domainempresa.Empresa,
	criarAdministrador func(empresaID int64) (*domainusuario.Usuario, error),
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}

		u, err := criarAdministrador(e.ID)
		if err != nil {
			return err
		}

		return tx.Create(u).Error
	})
}
