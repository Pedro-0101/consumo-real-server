package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/frentista"
)

type FrentistaGORMRepository struct {
	db *gorm.DB
}

func NewFrentistaGORMRepository(db *gorm.DB) *FrentistaGORMRepository {
	return &FrentistaGORMRepository{db: db}
}

func (r *FrentistaGORMRepository) Create(ctx context.Context, f *frentista.Frentista) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *FrentistaGORMRepository) Update(ctx context.Context, f *frentista.Frentista) error {
	return r.db.WithContext(ctx).Save(f).Error
}

func (r *FrentistaGORMRepository) FindByID(ctx context.Context, id int64) (*frentista.Frentista, error) {
	var f frentista.Frentista
	if err := r.db.WithContext(ctx).First(&f, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, frentista.ErrNaoEncontrado
		}
		return nil, err
	}
	return &f, nil
}

func (r *FrentistaGORMRepository) List(ctx context.Context, filter frentista.ListFilter) ([]frentista.Frentista, error) {
	q := r.db.WithContext(ctx).Model(&frentista.Frentista{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.UsuarioID > 0 {
		q = q.Where("usuario_id = ?", filter.UsuarioID)
	}
	if filter.Matricula != "" {
		q = q.Where("matricula = ?", filter.Matricula)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []frentista.Frentista
	if err := q.Order("nome asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
