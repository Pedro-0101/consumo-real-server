package seeds

import (
	"strings"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/usuario"
)

type AdminBaseGORMRepository struct {
	db *gorm.DB
}

func NewAdminBaseGORMRepository(db *gorm.DB) *AdminBaseGORMRepository {
	return &AdminBaseGORMRepository{db: db}
}

func (r *AdminBaseGORMRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&usuario.Usuario{}).
		Where("lower(email) = ?", strings.ToLower(email)).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AdminBaseGORMRepository) Save(u *usuario.Usuario) error {
	return r.db.Create(u).Error
}
