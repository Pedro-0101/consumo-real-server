package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/bomba"
)

type BombaGORMRepository struct {
	db *gorm.DB
}

func NewBombaGORMRepository(db *gorm.DB) *BombaGORMRepository {
	return &BombaGORMRepository{db: db}
}

func (r *BombaGORMRepository) Create(ctx context.Context, b *bomba.Bomba) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *BombaGORMRepository) Update(ctx context.Context, b *bomba.Bomba) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *BombaGORMRepository) FindByID(ctx context.Context, id int64) (*bomba.Bomba, error) {
	var b bomba.Bomba
	if err := r.db.WithContext(ctx).Preload("Bicos").First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bomba.ErrNaoEncontrado
		}
		return nil, err
	}
	return &b, nil
}

func (r *BombaGORMRepository) List(ctx context.Context, filter bomba.ListFilter) ([]bomba.Bomba, error) {
	q := r.db.WithContext(ctx).Model(&bomba.Bomba{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.LocalID > 0 {
		q = q.Where("local_id = ?", filter.LocalID)
	}
	if filter.ReservatorioID > 0 {
		q = q.Where("reservatorio_id = ?", filter.ReservatorioID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []bomba.Bomba
	if err := q.Preload("Bicos").Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *BombaGORMRepository) AdicionarBico(ctx context.Context, bico *bomba.Bico) error {
	return r.db.WithContext(ctx).Create(bico).Error
}

func (r *BombaGORMRepository) DesativarBico(ctx context.Context, bombaID, bicoID int64) error {
	result := r.db.WithContext(ctx).Model(&bomba.Bico{}).
		Where("id = ? AND bomba_id = ?", bicoID, bombaID).
		Update("ativo", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return bomba.ErrBicoNaoEncontrado
	}
	return nil
}
