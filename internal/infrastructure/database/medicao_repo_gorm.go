package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/medicao"
)

type MedicaoGORMRepository struct {
	db *gorm.DB
}

func NewMedicaoGORMRepository(db *gorm.DB) *MedicaoGORMRepository {
	return &MedicaoGORMRepository{db: db}
}

func (r *MedicaoGORMRepository) Create(ctx context.Context, m *medicao.Medicao) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MedicaoGORMRepository) Update(ctx context.Context, m *medicao.Medicao) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MedicaoGORMRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&medicao.Medicao{}, id).Error
}

func (r *MedicaoGORMRepository) FindByID(ctx context.Context, id int64) (*medicao.Medicao, error) {
	var m medicao.Medicao
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, medicao.ErrNaoEncontrado
		}
		return nil, err
	}
	return &m, nil
}

func (r *MedicaoGORMRepository) List(ctx context.Context, filter medicao.ListFilter) ([]medicao.Medicao, error) {
	q := r.db.WithContext(ctx).Model(&medicao.Medicao{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.ReservatorioID > 0 {
		q = q.Where("reservatorio_id = ?", filter.ReservatorioID)
	}

	var list []medicao.Medicao
	if err := q.Order("data desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
