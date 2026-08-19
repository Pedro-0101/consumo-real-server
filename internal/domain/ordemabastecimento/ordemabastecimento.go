package ordemabastecimento

import (
	"errors"
	"strings"
	"time"

	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/timeutil"
)

type Status string

const (
	StatusAberta     Status = "ABERTA"
	StatusAutorizada Status = "AUTORIZADA"
	StatusConcluida  Status = "CONCLUIDA"
	StatusCancelada  Status = "CANCELADA"
)

var (
	ErrEmpresaObrigatorio    = errors.New("empresa é obrigatória")
	ErrPatrimonioObrigatorio = errors.New("patrimônio é obrigatório")
	ErrNumeroObrigatorio     = errors.New("número é obrigatório")
	ErrQuantidadeInvalida    = errors.New("quantidade deve ser maior que zero")
	ErrStatusInvalido        = errors.New("operação não permitida para o status atual")
	ErrOrdemVencida          = errors.New("ordem de abastecimento vencida")
	ErrQuantidadeExcedida    = errors.New("quantidade abastecida excede a autorizada")
)

type OrdemAbastecimento struct {
	ID                   int64      `gorm:"primaryKey" json:"id"`
	EmpresaID            int64      `gorm:"not null;uniqueIndex:uk_empresa_numero,priority:1" json:"empresaID"`
	Numero               string     `gorm:"size:50;not null;uniqueIndex:uk_empresa_numero,priority:2" json:"numero"`
	PatrimonioID         int64      `gorm:"not null;index" json:"patrimonioID"`
	QuantidadeAutorizada float64    `gorm:"not null" json:"quantidadeAutorizada"`
	QuantidadeAbastecida float64    `gorm:"not null;default:0" json:"quantidadeAbastecida"`
	Status               Status     `gorm:"type:varchar(20);not null;index" json:"status"`
	DataEmissao          time.Time  `gorm:"not null" json:"dataEmissao"`
	DataValidade         *time.Time `gorm:"index" json:"dataValidade"`
	shared.AuditFields   `gorm:"embedded;embeddedPrefix:"`
}

func NewOrdemAbastecimento(empresaID, patrimonioID int64, numero string, quantidadeAutorizada float64, dataValidade *time.Time) (*OrdemAbastecimento, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatorio
	}
	if patrimonioID <= 0 {
		return nil, ErrPatrimonioObrigatorio
	}
	if strings.TrimSpace(numero) == "" {
		return nil, ErrNumeroObrigatorio
	}
	if quantidadeAutorizada <= 0 {
		return nil, ErrQuantidadeInvalida
	}

	return &OrdemAbastecimento{
		EmpresaID:            empresaID,
		Numero:               strings.TrimSpace(numero),
		PatrimonioID:         patrimonioID,
		QuantidadeAutorizada: quantidadeAutorizada,
		Status:               StatusAberta,
		DataEmissao:          timeutil.Now(),
		DataValidade:         dataValidade,
	}, nil
}

func (o *OrdemAbastecimento) Autorizar() error {
	if o.Status != StatusAberta {
		return ErrStatusInvalido
	}
	o.Status = StatusAutorizada
	return nil
}

func (o *OrdemAbastecimento) RegistrarAbastecimento(quantidade float64) error {
	if quantidade <= 0 {
		return ErrQuantidadeInvalida
	}
	if o.Status != StatusAutorizada {
		return ErrStatusInvalido
	}
	if o.DataValidade != nil && timeutil.Now().After(*o.DataValidade) {
		return ErrOrdemVencida
	}
	if o.QuantidadeAbastecida+quantidade > o.QuantidadeAutorizada {
		return ErrQuantidadeExcedida
	}
	o.QuantidadeAbastecida += quantidade
	if o.QuantidadeAbastecida >= o.QuantidadeAutorizada {
		o.Status = StatusConcluida
	}
	return nil
}

func (o *OrdemAbastecimento) Cancelar() error {
	if o.Status == StatusConcluida {
		return ErrStatusInvalido
	}
	o.Status = StatusCancelada
	return nil
}
