package combustivel

import (
	"errors"

	"consumo-real-server/internal/shared"
)

type Tipo string

const (
	TipoGasolina Tipo = "GASOLINA"
	TipoEtanol   Tipo = "ETANOL"
	TipoDiesel   Tipo = "DIESEL"
	TipoGNV      Tipo = "GNV"
)

type Unidade string

const (
	UnidadeLitro Unidade = "LITRO"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrTipoInvalido       = errors.New("tipo de combustível inválido")
	ErrUnidadeInvalida    = errors.New("unidade de medida inválida")
	ErrPrecoInvalido      = errors.New("preço não pode ser negativo")
	ErrDensidadeInvalida  = errors.New("densidade deve ser maior que zero")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
)

type Combustivel struct {
	ID                 int64   `gorm:"primaryKey" json:"id"`
	EmpresaID          int64   `gorm:"not null;index" json:"empresaID"`
	Nome               string  `gorm:"size:255;not null" json:"nome"`
	Tipo               Tipo    `gorm:"type:varchar(30);not null" json:"tipo"`
	Unidade            Unidade `gorm:"type:varchar(20);not null" json:"unidade"`
	Densidade          float64 `gorm:"not null" json:"densidade"`
	PrecoCusto         float64 `gorm:"not null;default:0" json:"precoCusto"`
	PrecoVenda         float64 `gorm:"not null;default:0" json:"precoVenda"`
	Ativo              bool    `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewCombustivel(empresaID int64, nome string, tipo Tipo, unidade Unidade, densidade, precoCusto, precoVenda float64) (*Combustivel, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if nome == "" {
		return nil, ErrNomeObrigatorio
	}
	if !tipo.isValid() {
		return nil, ErrTipoInvalido
	}
	if !unidade.isValid() {
		return nil, ErrUnidadeInvalida
	}
	if densidade <= 0 {
		return nil, ErrDensidadeInvalida
	}
	if precoCusto < 0 || precoVenda < 0 {
		return nil, ErrPrecoInvalido
	}

	return &Combustivel{
		EmpresaID:  empresaID,
		Nome:       nome,
		Tipo:       tipo,
		Unidade:    unidade,
		Densidade:  densidade,
		PrecoCusto: precoCusto,
		PrecoVenda: precoVenda,
		Ativo:      true,
	}, nil
}

func (c *Combustivel) Atualizar(nome string, tipo Tipo, unidade Unidade, densidade float64) error {
	if nome == "" {
		return ErrNomeObrigatorio
	}
	if !tipo.isValid() {
		return ErrTipoInvalido
	}
	if !unidade.isValid() {
		return ErrUnidadeInvalida
	}
	if densidade <= 0 {
		return ErrDensidadeInvalida
	}

	c.Nome = nome
	c.Tipo = tipo
	c.Unidade = unidade
	c.Densidade = densidade
	return nil
}

func (c *Combustivel) AtualizarPrecos(precoCusto, precoVenda float64) error {
	if precoCusto < 0 || precoVenda < 0 {
		return ErrPrecoInvalido
	}
	c.PrecoCusto = precoCusto
	c.PrecoVenda = precoVenda
	return nil
}

func (c *Combustivel) Desativar() {
	c.Ativo = false
}

func (t Tipo) isValid() bool {
	switch t {
	case TipoGasolina, TipoEtanol, TipoDiesel, TipoGNV:
		return true
	}
	return false
}

func (u Unidade) isValid() bool {
	switch u {
	case UnidadeLitro:
		return true
	}
	return false
}
