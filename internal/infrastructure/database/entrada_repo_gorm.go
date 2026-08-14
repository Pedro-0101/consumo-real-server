package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/entrada"
)

type EntradaGORMRepository struct {
	db *gorm.DB
}

func NewEntradaGORMRepository(db *gorm.DB) *EntradaGORMRepository {
	return &EntradaGORMRepository{db: db}
}

func (r *EntradaGORMRepository) Create(ctx context.Context, e *entrada.Entrada) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *EntradaGORMRepository) Update(ctx context.Context, e *entrada.Entrada) error {
	return r.db.WithContext(ctx).Save(e).Error
}

func (r *EntradaGORMRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entrada.Entrada{}, id).Error
}

func (r *EntradaGORMRepository) FindByID(ctx context.Context, id int64) (*entrada.Entrada, error) {
	var e entrada.Entrada
	if err := r.db.WithContext(ctx).Preload("Combustivel").First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entrada.ErrNaoEncontrado
		}
		return nil, err
	}
	return &e, nil
}

func (r *EntradaGORMRepository) List(ctx context.Context, filter entrada.ListFilter) ([]entrada.Entrada, error) {
	q := r.db.WithContext(ctx).Model(&entrada.Entrada{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.FornecedorID > 0 {
		q = q.Where("fornecedor_id = ?", filter.FornecedorID)
	}
	if filter.ReservatorioID > 0 {
		q = q.Where("reservatorio_id = ?", filter.ReservatorioID)
	}
	if filter.CombustivelID > 0 {
		q = q.Where("combustivel_id = ?", filter.CombustivelID)
	}

	var list []entrada.Entrada
	if err := q.Preload("Combustivel").Order("data desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
