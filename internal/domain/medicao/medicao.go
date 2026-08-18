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
	ID                 int64     `gorm:"primaryKey" json:"id"`
	EmpresaID          int64     `gorm:"not null;index" json:"empresaID"`
	ReservatorioID     int64     `gorm:"not null;index" json:"reservatorioID"`
	NivelCalculado     float64   `gorm:"not null" json:"nivelCalculado"`
	NivelMedido        float64   `gorm:"not null" json:"nivelMedido"`
	Diferenca          float64   `gorm:"not null" json:"diferenca"`
	Data               time.Time `gorm:"not null;index" json:"data"`
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
