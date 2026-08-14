package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/preco"
)

type PrecoGORMRepository struct {
	db *gorm.DB
}

func NewPrecoGORMRepository(db *gorm.DB) *PrecoGORMRepository {
	return &PrecoGORMRepository{db: db}
}

func (r *PrecoGORMRepository) Create(ctx context.Context, p *preco.Preco) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PrecoGORMRepository) Update(ctx context.Context, p *preco.Preco) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *PrecoGORMRepository) FindByID(ctx context.Context, id int64) (*preco.Preco, error) {
	var p preco.Preco
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, preco.ErrNaoEncontrado
		}
		return nil, err
	}
	return &p, nil
}

func (r *PrecoGORMRepository) List(ctx context.Context, filter preco.ListFilter) ([]preco.Preco, error) {
	q := r.db.WithContext(ctx).Model(&preco.Preco{})
	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.CombustivelID > 0 {
		q = q.Where("combustivel_id = ?", filter.CombustivelID)
	}
	if filter.Ativo != nil {
		q = q.Where("ativo = ?", *filter.Ativo)
	}

	var list []preco.Preco
	if err := q.Order("vigencia_inicio desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
