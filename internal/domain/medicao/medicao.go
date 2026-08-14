package medicao

import (
	"errors"
	"time"

	"consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/timeutil"
)

var (
	ErrEmpresaObrigatoria   = errors.New("empresa é obrigatória")
	ErrReservatorioInvalido = errors.New("reservatório é obrigatório")
	ErrEmpresaIncompativel  = errors.New("reservatório pertence a outra empresa")
	ErrNivelMedidoInvalido  = errors.New("nível medido não pode ser negativo")
)

type Medicao struct {
	ID                 int64     `gorm:"primaryKey"`
	EmpresaID          int64     `gorm:"not null;index"`
	ReservatorioID     int64     `gorm:"not null;index"`
	NivelCalculado     float64   `gorm:"not null"`
	NivelMedido        float64   `gorm:"not null"`
	Diferenca          float64   `gorm:"not null"`
	Data               time.Time `gorm:"not null;index"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewMedicao(empresaID int64, reservatorio *reservatorio.Reservatorio, nivelMedido float64) (*Medicao, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if reservatorio == nil || reservatorio.ID == 0 {
		return nil, ErrReservatorioInvalido
	}
	if reservatorio.EmpresaID != empresaID {
		return nil, ErrEmpresaIncompativel
	}
	if nivelMedido < 0 {
		return nil, ErrNivelMedidoInvalido
	}

	nivelCalculado := reservatorio.NivelAtual
	if err := reservatorio.CorrigirNivel(nivelMedido); err != nil {
		return nil, err
	}

	return &Medicao{
		EmpresaID:      empresaID,
		ReservatorioID: reservatorio.ID,
		NivelCalculado: nivelCalculado,
		NivelMedido:    nivelMedido,
		Diferenca:      nivelMedido - nivelCalculado,
		Data:           timeutil.Now(),
	}, nil
}
