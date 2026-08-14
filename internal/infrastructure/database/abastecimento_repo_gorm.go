package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/abastecimento"
)

type AbastecimentoGORMRepository struct {
	db *gorm.DB
}

func NewAbastecimentoGORMRepository(db *gorm.DB) *AbastecimentoGORMRepository {
	return &AbastecimentoGORMRepository{db: db}
}

func (r *AbastecimentoGORMRepository) Create(ctx context.Context, a *abastecimento.Abastecimento) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AbastecimentoGORMRepository) Update(ctx context.Context, a *abastecimento.Abastecimento) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AbastecimentoGORMRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&abastecimento.Abastecimento{}, id).Error
}

func (r *AbastecimentoGORMRepository) FindByID(ctx context.Context, id int64) (*abastecimento.Abastecimento, error) {
	var a abastecimento.Abastecimento
	if err := r.db.WithContext(ctx).Preload("Combustivel").First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, abastecimento.ErrNaoEncontrado
		}
		return nil, err
	}
	return &a, nil
}

func (r *AbastecimentoGORMRepository) List(ctx context.Context, filter abastecimento.ListFilter) ([]abastecimento.Abastecimento, error) {
	q := r.db.WithContext(ctx).Model(&abastecimento.Abastecimento{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.LocalID > 0 {
		q = q.Where("local_id = ?", filter.LocalID)
	}
	if filter.BombaID > 0 {
		q = q.Where("bomba_id = ?", filter.BombaID)
	}
	if filter.PatrimonioID > 0 {
		q = q.Where("patrimonio_id = ?", filter.PatrimonioID)
	}
	if filter.FrentistaID > 0 {
		q = q.Where("frentista_id = ?", filter.FrentistaID)
	}
	if filter.CombustivelID > 0 {
		q = q.Where("combustivel_id = ?", filter.CombustivelID)
	}
	if filter.Tipo != "" {
		q = q.Where("tipo = ?", filter.Tipo)
	}

	var list []abastecimento.Abastecimento
	if err := q.Preload("Combustivel").Order("data desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
