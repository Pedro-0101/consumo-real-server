package reservatorio

import (
	"errors"

	"consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio      = errors.New("nome é obrigatório")
	ErrCapacidadeInvalida   = errors.New("capacidade deve ser maior que zero")
	ErrNivelInicialInvalido = errors.New("nível inicial deve estar entre zero e a capacidade")
	ErrNivelMinimoInvalido  = errors.New("nível mínimo deve estar entre zero e a capacidade")
	ErrCombustivelInvalido  = errors.New("combustível é obrigatório")
	ErrEmpresaObrigatoria   = errors.New("empresa é obrigatória")
	ErrEmpresaIncompativel  = errors.New("combustível pertence a outra empresa")
	ErrCapacidadeExcedida   = errors.New("entrada excede a capacidade do reservatório")
	ErrNivelInsuficiente    = errors.New("saída maior que o nível atual de combustível")
	ErrQuantidadeInvalida   = errors.New("quantidade deve ser maior que zero")
	ErrReservatorioInativo  = errors.New("reservatório inativo")
	ErrNivelMedidoInvalido  = errors.New("nível medido não pode ser negativo")
)

type Reservatorio struct {
	ID                 int64                   `gorm:"primaryKey" json:"id"`
	EmpresaID          int64                   `gorm:"not null;index" json:"empresaID"`
	Nome               string                  `gorm:"size:255;not null" json:"nome"`
	Capacidade         float64                 `gorm:"not null" json:"capacidade"`
	NivelAtual         float64                 `gorm:"not null;default:0" json:"nivelAtual"`
	NivelMinimo        float64                 `gorm:"not null;default:0" json:"nivelMinimo"`
	CombustivelID      int64                   `gorm:"not null;index" json:"combustivelID"`
	Combustivel        combustivel.Combustivel `gorm:"foreignKey:CombustivelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"combustivel"`
	Ativo              bool                    `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewReservatorio(empresaID int64, nome string, capacidade, nivelInicial, nivelMinimo float64, combustivel combustivel.Combustivel) (*Reservatorio, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if nome == "" {
		return nil, ErrNomeObrigatorio
	}
	if capacidade <= 0 {
		return nil, ErrCapacidadeInvalida
	}
	if nivelInicial < 0 || nivelInicial > capacidade {
		return nil, ErrNivelInicialInvalido
	}
	if nivelMinimo < 0 || nivelMinimo > capacidade {
		return nil, ErrNivelMinimoInvalido
	}
	if combustivel.ID == 0 {
		return nil, ErrCombustivelInvalido
	}
	if combustivel.EmpresaID != empresaID {
		return nil, ErrEmpresaIncompativel
	}

	return &Reservatorio{
		EmpresaID:     empresaID,
		Nome:          nome,
		Capacidade:    capacidade,
		NivelAtual:    nivelInicial,
		NivelMinimo:   nivelMinimo,
		CombustivelID: combustivel.ID,
		Combustivel:   combustivel,
		Ativo:         true,
	}, nil
}

func (r *Reservatorio) Entrada(quantidade float64) error {
	if !r.Ativo {
		return ErrReservatorioInativo
	}
	if quantidade <= 0 {
		return ErrQuantidadeInvalida
	}
	if r.NivelAtual+quantidade > r.Capacidade {
		return ErrCapacidadeExcedida
	}
	r.NivelAtual += quantidade
	return nil
}

func (r *Reservatorio) Saida(quantidade float64) error {
	if !r.Ativo {
		return ErrReservatorioInativo
	}
	if quantidade <= 0 {
		return ErrQuantidadeInvalida
	}
	if quantidade > r.NivelAtual {
		return ErrNivelInsuficiente
	}
	r.NivelAtual -= quantidade
	return nil
}

func (r *Reservatorio) NivelPercentual() float64 {
	if r.Capacidade == 0 {
		return 0
	}
	return (r.NivelAtual / r.Capacidade) * 100
}

func (r *Reservatorio) CorrigirNivel(nivelMedido float64) error {
	if nivelMedido < 0 {
		return ErrNivelMedidoInvalido
	}
	if nivelMedido > r.Capacidade {
		return ErrCapacidadeExcedida
	}
	r.NivelAtual = nivelMedido
	return nil
}

func (r *Reservatorio) AbaixoDoMinimo() bool {
	return r.NivelAtual < r.NivelMinimo
}

func (r *Reservatorio) Desativar() {
	r.Ativo = false
}

func (r *Reservatorio) Atualizar(nome string, capacidade, nivelMinimo float64, combustivel combustivel.Combustivel) error {
	if nome == "" {
		return ErrNomeObrigatorio
	}
	if capacidade <= 0 {
		return ErrCapacidadeInvalida
	}
	if nivelMinimo < 0 || nivelMinimo > capacidade {
		return ErrNivelMinimoInvalido
	}
	if combustivel.ID == 0 {
		return ErrCombustivelInvalido
	}
	if combustivel.EmpresaID != r.EmpresaID {
		return ErrEmpresaIncompativel
	}

	r.Nome = nome
	r.Capacidade = capacidade
	r.NivelMinimo = nivelMinimo
	r.CombustivelID = combustivel.ID
	r.Combustivel = combustivel
	return nil
}
