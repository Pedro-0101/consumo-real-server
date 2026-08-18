package preco

import (
	"errors"
	"time"

	"consumo-real-server/internal/shared"
)

var (
	ErrEmpresaObrigatoria     = errors.New("empresa é obrigatória")
	ErrCombustivelObrigatorio = errors.New("combustível é obrigatório")
	ErrPrecoInvalido          = errors.New("preço não pode ser negativo")
	ErrVigenciaInvalida       = errors.New("vigência inválida")
)

type Preco struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	EmpresaID          int64      `gorm:"not null;index" json:"empresaID"`
	CombustivelID      int64      `gorm:"not null;index" json:"combustivelID"`
	PrecoCusto         float64    `gorm:"not null" json:"precoCusto"`
	PrecoVenda         float64    `gorm:"not null" json:"precoVenda"`
	VigenciaInicio     time.Time  `gorm:"not null" json:"vigenciaInicio"`
	VigenciaFim        *time.Time `gorm:"index" json:"vigenciaFim"`
	Ativo              bool       `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewPreco(empresaID, combustivelID int64, precoCusto, precoVenda float64, vigenciaInicio time.Time) (*Preco, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if combustivelID <= 0 {
		return nil, ErrCombustivelObrigatorio
	}
	if precoCusto < 0 || precoVenda < 0 {
		return nil, ErrPrecoInvalido
	}
	if vigenciaInicio.IsZero() {
		return nil, ErrVigenciaInvalida
	}

	return &Preco{
		EmpresaID:      empresaID,
		CombustivelID:  combustivelID,
		PrecoCusto:     precoCusto,
		PrecoVenda:     precoVenda,
		VigenciaInicio: vigenciaInicio,
		Ativo:          true,
	}, nil
}

func (p *Preco) Encerrar(dataFim time.Time) error {
	if dataFim.Before(p.VigenciaInicio) {
		return ErrVigenciaInvalida
	}
	p.VigenciaFim = &dataFim
	p.Ativo = false
	return nil
}

func (p *Preco) Vigente(data time.Time) bool {
	if !p.Ativo {
		return false
	}
	if data.Before(p.VigenciaInicio) {
		return false
	}
	if p.VigenciaFim != nil && data.After(*p.VigenciaFim) {
		return false
	}
	return true
}
