package database

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/usuario"
)

type UsuarioGORMRepository struct {
	db *gorm.DB
}

func NewUsuarioGORMRepository(db *gorm.DB) *UsuarioGORMRepository {
	return &UsuarioGORMRepository{db: db}
}

func (r *UsuarioGORMRepository) Create(ctx context.Context, u *usuario.Usuario) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UsuarioGORMRepository) Update(ctx context.Context, u *usuario.Usuario) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UsuarioGORMRepository) FindByID(ctx context.Context, id int64) (*usuario.Usuario, error) {
	var u usuario.Usuario
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usuario.ErrNaoEncontrado
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioGORMRepository) FindByEmail(ctx context.Context, email string) (*usuario.Usuario, error) {
	var u usuario.Usuario
	if err := r.db.WithContext(ctx).
		Where("lower(email) = ?", strings.ToLower(email)).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usuario.ErrNaoEncontrado
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioGORMRepository) List(ctx context.Context, filter usuario.ListFilter) ([]usuario.Usuario, error) {
	q := r.db.WithContext(ctx).Model(&usuario.Usuario{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.Papel != "" {
		q = q.Where("papel = ?", filter.Papel)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []usuario.Usuario
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
