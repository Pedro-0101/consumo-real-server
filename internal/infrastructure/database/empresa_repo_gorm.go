package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/empresa"
)

type EmpresaGORMRepository struct {
	db *gorm.DB
}

func NewEmpresaGORMRepository(db *gorm.DB) *EmpresaGORMRepository {
	return &EmpresaGORMRepository{db: db}
}

func (r *EmpresaGORMRepository) Create(ctx context.Context, e *empresa.Empresa) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *EmpresaGORMRepository) Update(ctx context.Context, e *empresa.Empresa) error {
	return r.db.WithContext(ctx).Save(e).Error
}

func (r *EmpresaGORMRepository) FindByID(ctx context.Context, id int64) (*empresa.Empresa, error) {
	var e empresa.Empresa
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, empresa.ErrNaoEncontrado
		}
		return nil, err
	}
	return &e, nil
}

func (r *EmpresaGORMRepository) List(ctx context.Context, filter empresa.ListFilter) ([]empresa.Empresa, error) {
	q := r.db.WithContext(ctx).Model(&empresa.Empresa{})
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []empresa.Empresa
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
