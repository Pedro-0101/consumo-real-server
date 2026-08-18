package database

import (
	"context"

	"gorm.io/gorm"

	"consumo-real-server/internal/domain/auditoria"
)

type AuditoriaGORMRepository struct {
	db *gorm.DB
}

func NewAuditoriaGORMRepository(db *gorm.DB) *AuditoriaGORMRepository {
	return &AuditoriaGORMRepository{db: db}
}

func (r *AuditoriaGORMRepository) Create(ctx context.Context, a *auditoria.Auditoria) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AuditoriaGORMRepository) List(ctx context.Context, filter auditoria.ListFilter) ([]auditoria.Auditoria, error) {
	q := r.db.WithContext(ctx).Model(&auditoria.Auditoria{})

	if filter.EmpresaID > 0 {
		q = q.Where("empresa_id = ?", filter.EmpresaID)
	}
	if filter.Entidade != "" {
		q = q.Where("entidade = ?", filter.Entidade)
	}
	if filter.Operacao != "" {
		q = q.Where("operacao = ?", filter.Operacao)
	}
	if filter.UsuarioID > 0 {
		q = q.Where("usuario_id = ?", filter.UsuarioID)
	}
	if filter.De != nil {
		q = q.Where("created_at >= ?", *filter.De)
	}
	if filter.Ate != nil {
		q = q.Where("created_at <= ?", *filter.Ate)
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	}

	var list []auditoria.Auditoria
	if err := q.Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}