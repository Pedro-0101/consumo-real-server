package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/combustivel"
)

type CombustivelGORMRepository struct {
	db *gorm.DB
}

func NewCombustivelGORMRepository(db *gorm.DB) *CombustivelGORMRepository {
	return &CombustivelGORMRepository{db: db}
}

func (r *CombustivelGORMRepository) Create(ctx context.Context, c *combustivel.Combustivel) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CombustivelGORMRepository) Update(ctx context.Context, c *combustivel.Combustivel) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CombustivelGORMRepository) FindByID(ctx context.Context, id int64) (*combustivel.Combustivel, error) {
	var c combustivel.Combustivel
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, combustivel.ErrNaoEncontrado
		}
		return nil, err
	}
	return &c, nil
}

func (r *CombustivelGORMRepository) List(ctx context.Context, filter combustivel.ListFilter) ([]combustivel.Combustivel, error) {
	q := r.db.WithContext(ctx).Model(&combustivel.Combustivel{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []combustivel.Combustivel
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
