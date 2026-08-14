package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/patrimonio"
)

type PatrimonioGORMRepository struct {
	db *gorm.DB
}

func NewPatrimonioGORMRepository(db *gorm.DB) *PatrimonioGORMRepository {
	return &PatrimonioGORMRepository{db: db}
}

func (r *PatrimonioGORMRepository) Create(ctx context.Context, p *patrimonio.Patrimonio) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PatrimonioGORMRepository) Update(ctx context.Context, p *patrimonio.Patrimonio) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *PatrimonioGORMRepository) FindByID(ctx context.Context, id int64) (*patrimonio.Patrimonio, error) {
	var p patrimonio.Patrimonio
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, patrimonio.ErrNaoEncontrado
		}
		return nil, err
	}
	return &p, nil
}

func (r *PatrimonioGORMRepository) List(ctx context.Context, filter patrimonio.ListFilter) ([]patrimonio.Patrimonio, error) {
	q := r.db.WithContext(ctx).Model(&patrimonio.Patrimonio{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.UnidadeAdministrativaID > 0 {
		q = q.Where("unidade_administrativa_id = ?", filter.UnidadeAdministrativaID)
	}
	if filter.Tipo != "" {
		q = q.Where("tipo = ?", filter.Tipo)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []patrimonio.Patrimonio
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
