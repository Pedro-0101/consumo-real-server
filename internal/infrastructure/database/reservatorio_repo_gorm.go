package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/reservatorio"
)

type ReservatorioGORMRepository struct {
	db *gorm.DB
}

func NewReservatorioGORMRepository(db *gorm.DB) *ReservatorioGORMRepository {
	return &ReservatorioGORMRepository{db: db}
}

func (r *ReservatorioGORMRepository) Create(ctx context.Context, reserv *reservatorio.Reservatorio) error {
	return r.db.WithContext(ctx).Create(reserv).Error
}

func (r *ReservatorioGORMRepository) Update(ctx context.Context, reserv *reservatorio.Reservatorio) error {
	return r.db.WithContext(ctx).Save(reserv).Error
}

func (r *ReservatorioGORMRepository) FindByID(ctx context.Context, id int64) (*reservatorio.Reservatorio, error) {
	var reserv reservatorio.Reservatorio
	if err := r.db.WithContext(ctx).Preload("Combustivel").First(&reserv, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, reservatorio.ErrNaoEncontrado
		}
		return nil, err
	}
	return &reserv, nil
}

func (r *ReservatorioGORMRepository) List(ctx context.Context, filter reservatorio.ListFilter) ([]reservatorio.Reservatorio, error) {
	q := r.db.WithContext(ctx).Model(&reservatorio.Reservatorio{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.CombustivelID > 0 {
		q = q.Where("combustivel_id = ?", filter.CombustivelID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []reservatorio.Reservatorio
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
