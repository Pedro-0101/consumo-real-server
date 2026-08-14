package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/fornecedor"
)

type FornecedorGORMRepository struct {
	db *gorm.DB
}

func NewFornecedorGORMRepository(db *gorm.DB) *FornecedorGORMRepository {
	return &FornecedorGORMRepository{db: db}
}

func (r *FornecedorGORMRepository) Create(ctx context.Context, f *fornecedor.Fornecedor) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *FornecedorGORMRepository) Update(ctx context.Context, f *fornecedor.Fornecedor) error {
	return r.db.WithContext(ctx).Save(f).Error
}

func (r *FornecedorGORMRepository) FindByID(ctx context.Context, id int64) (*fornecedor.Fornecedor, error) {
	var f fornecedor.Fornecedor
	if err := r.db.WithContext(ctx).First(&f, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fornecedor.ErrNaoEncontrado
		}
		return nil, err
	}
	return &f, nil
}

func (r *FornecedorGORMRepository) List(ctx context.Context, filter fornecedor.ListFilter) ([]fornecedor.Fornecedor, error) {
	q := r.db.WithContext(ctx).Model(&fornecedor.Fornecedor{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.CNPJ != "" {
		q = q.Where("cnpj = ?", filter.CNPJ)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []fornecedor.Fornecedor
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
