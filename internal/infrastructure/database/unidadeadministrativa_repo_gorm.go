package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/unidadeadministrativa"
)

type UnidadeAdministrativaGORMRepository struct {
	db *gorm.DB
}

func NewUnidadeAdministrativaGORMRepository(db *gorm.DB) *UnidadeAdministrativaGORMRepository {
	return &UnidadeAdministrativaGORMRepository{db: db}
}

func (r *UnidadeAdministrativaGORMRepository) Create(ctx context.Context, u *unidadeadministrativa.UnidadeAdministrativa) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UnidadeAdministrativaGORMRepository) Update(ctx context.Context, u *unidadeadministrativa.UnidadeAdministrativa) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UnidadeAdministrativaGORMRepository) FindByID(ctx context.Context, id int64) (*unidadeadministrativa.UnidadeAdministrativa, error) {
	var u unidadeadministrativa.UnidadeAdministrativa
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, unidadeadministrativa.ErrNaoEncontrado
		}
		return nil, err
	}
	return &u, nil
}

func (r *UnidadeAdministrativaGORMRepository) List(ctx context.Context, filter unidadeadministrativa.ListFilter) ([]unidadeadministrativa.UnidadeAdministrativa, error) {
	q := r.db.WithContext(ctx).Model(&unidadeadministrativa.UnidadeAdministrativa{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []unidadeadministrativa.UnidadeAdministrativa
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
