package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/ordemabastecimento"
)

type OrdemAbastecimentoGORMRepository struct {
	db *gorm.DB
}

func NewOrdemAbastecimentoGORMRepository(db *gorm.DB) *OrdemAbastecimentoGORMRepository {
	return &OrdemAbastecimentoGORMRepository{db: db}
}

func (r *OrdemAbastecimentoGORMRepository) Create(ctx context.Context, o *ordemabastecimento.OrdemAbastecimento) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *OrdemAbastecimentoGORMRepository) Update(ctx context.Context, o *ordemabastecimento.OrdemAbastecimento) error {
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *OrdemAbastecimentoGORMRepository) FindByID(ctx context.Context, id int64) (*ordemabastecimento.OrdemAbastecimento, error) {
	var o ordemabastecimento.OrdemAbastecimento
	if err := r.db.WithContext(ctx).First(&o, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ordemabastecimento.ErrNaoEncontrado
		}
		return nil, err
	}
	return &o, nil
}

func (r *OrdemAbastecimentoGORMRepository) List(ctx context.Context, filter ordemabastecimento.ListFilter) ([]ordemabastecimento.OrdemAbastecimento, error) {
	q := r.db.WithContext(ctx).Model(&ordemabastecimento.OrdemAbastecimento{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.PatrimonioID > 0 {
		q = q.Where("patrimonio_id = ?", filter.PatrimonioID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	var list []ordemabastecimento.OrdemAbastecimento
	if err := q.Order("data_emissao desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
