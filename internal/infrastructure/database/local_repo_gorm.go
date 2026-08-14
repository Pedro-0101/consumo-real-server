package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/local"
)

type LocalGORMRepository struct {
	db *gorm.DB
}

func NewLocalGORMRepository(db *gorm.DB) *LocalGORMRepository {
	return &LocalGORMRepository{db: db}
}

func (r *LocalGORMRepository) Create(ctx context.Context, l *local.Local) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *LocalGORMRepository) Update(ctx context.Context, l *local.Local) error {
	return r.db.WithContext(ctx).Save(l).Error
}

func (r *LocalGORMRepository) FindByID(ctx context.Context, id int64) (*local.Local, error) {
	var l local.Local
	if err := r.db.WithContext(ctx).First(&l, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, local.ErrNaoEncontrado
		}
		return nil, err
	}
	return &l, nil
}

func (r *LocalGORMRepository) List(ctx context.Context, filter local.ListFilter) ([]local.Local, error) {
	q := r.db.WithContext(ctx).Model(&local.Local{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.UnidadeAdministrativaID > 0 {
		q = q.Where("unidade_administrativa_id = ?", filter.UnidadeAdministrativaID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []local.Local
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
