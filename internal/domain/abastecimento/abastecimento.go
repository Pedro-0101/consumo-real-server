package abastecimento

import (
	"errors"
	"time"

	"consumo-real-server/internal/domain/bomba"
	"consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/timeutil"
)

type Tipo string

const (
	TipoAbastecimento Tipo = "ABASTECIMENTO"
	TipoTransferencia Tipo = "TRANSFERENCIA"
)

var (
	ErrEmpresaObrigatoria          = errors.New("empresa é obrigatória")
	ErrLocalObrigatorio            = errors.New("local é obrigatório")
	ErrBombaObrigatoria            = errors.New("bomba é obrigatória")
	ErrBicoObrigatorio             = errors.New("bico é obrigatório")
	ErrBicoInvalido                = errors.New("bico não pertence à bomba")
	ErrPatrimonioObrigatorio       = errors.New("patrimônio é obrigatório")
	ErrFrentistaObrigatorio        = errors.New("frentista é obrigatório")
	ErrReservatorioInvalido        = errors.New("reservatório de origem é obrigatório")
	ErrReservatorioDestinoInvalido = errors.New("reservatório de destino é obrigatório")
	ErrReservatorioIgual           = errors.New("reservatório de origem e destino devem ser diferentes")
	ErrReservatorioIncompativel    = errors.New("bomba não está vinculada ao reservatório de origem")
	ErrEmpresaIncompativel         = errors.New("entidades pertencem a empresas diferentes")
	ErrCombustivelDiferente        = errors.New("reservatórios devem conter o mesmo combustível")
	ErrQuantidadeInvalida          = errors.New("quantidade deve ser maior que zero")
	ErrPrecoInvalido               = errors.New("preço unitário não pode ser negativo")
	ErrLeituraInvalida             = errors.New("leitura não pode ser negativa")
)

type Abastecimento struct {
	ID                    int64                   `gorm:"primaryKey" json:"id"`
	EmpresaID             int64                   `gorm:"not null;index" json:"empresaID"`
	LocalID               int64                   `gorm:"not null;index" json:"localID"`
	BombaID               int64                   `gorm:"not null;index" json:"bombaID"`
	BicoID                int64                   `gorm:"not null;index" json:"bicoID"`
	Tipo                  Tipo                    `gorm:"type:varchar(20);not null;index" json:"tipo"`
	Data                  time.Time               `gorm:"not null;index" json:"data"`
	Quantidade            float64                 `gorm:"not null" json:"quantidade"`
	PrecoUnitario         float64                 `gorm:"not null;default:0" json:"precoUnitario"`
	ValorTotal            float64                 `gorm:"not null;default:0" json:"valorTotal"`
	Odometro              float64                 `gorm:"default:0" json:"odometro"`
	Horimetro             float64                 `gorm:"default:0" json:"horimetro"`
	CombustivelID         int64                   `gorm:"not null;index" json:"combustivelID"`
	Combustivel           combustivel.Combustivel `gorm:"foreignKey:CombustivelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"combustivel"`
	PatrimonioID          int64                   `gorm:"not null;index" json:"patrimonioID"`
	FrentistaID           int64                   `gorm:"not null;index" json:"frentistaID"`
	ReservatorioOrigemID  int64                   `gorm:"not null;index" json:"reservatorioOrigemID"`
	ReservatorioDestinoID int64                   `gorm:"index" json:"reservatorioDestinoID"`
	shared.AuditFields    `gorm:"embedded;embeddedPrefix:"`
}

func NewAbastecimento(empresaID, localID int64, bomba *bomba.Bomba, origem *reservatorio.Reservatorio, frentistaID, patrimonioID, bicoID int64, quantidade, precoUnitario, odometro, horimetro float64) (*Abastecimento, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if localID <= 0 {
		return nil, ErrLocalObrigatorio
	}
	if bomba == nil || bomba.ID == 0 {
		return nil, ErrBombaObrigatoria
	}
	if origem == nil || origem.ID == 0 {
		return nil, ErrReservatorioInvalido
	}
	if bomba.ReservatorioID != origem.ID {
		return nil, ErrReservatorioIncompativel
	}
	if bomba.EmpresaID != empresaID || origem.EmpresaID != empresaID {
		return nil, ErrEmpresaIncompativel
	}
	if frentistaID <= 0 {
		return nil, ErrFrentistaObrigatorio
	}
	if patrimonioID <= 0 {
		return nil, ErrPatrimonioObrigatorio
	}
	if bicoID <= 0 {
		return nil, ErrBicoObrigatorio
	}
	if !bomba.TemBicoAtivo(bicoID) {
		return nil, ErrBicoInvalido
	}
	if quantidade <= 0 {
		return nil, ErrQuantidadeInvalida
	}
	if precoUnitario < 0 {
		return nil, ErrPrecoInvalido
	}
	if odometro < 0 || horimetro < 0 {
		return nil, ErrLeituraInvalida
	}
	if err := origem.Saida(quantidade); err != nil {
		return nil, err
	}

	return &Abastecimento{
		EmpresaID:            empresaID,
		LocalID:              localID,
		BombaID:              bomba.ID,
		BicoID:               bicoID,
		Tipo:                 TipoAbastecimento,
		Data:                 timeutil.Now(),
		Quantidade:           quantidade,
		PrecoUnitario:        precoUnitario,
		ValorTotal:           quantidade * precoUnitario,
		Odometro:             odometro,
		Horimetro:            horimetro,
		CombustivelID:        origem.Combustivel.ID,
		Combustivel:          origem.Combustivel,
		PatrimonioID:         patrimonioID,
		FrentistaID:          frentistaID,
		ReservatorioOrigemID: origem.ID,
	}, nil
}

func NewTransferencia(empresaID int64, origem, destino *reservatorio.Reservatorio, quantidade float64) (*Abastecimento, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if origem == nil || origem.ID == 0 {
		return nil, ErrReservatorioInvalido
	}
	if destino == nil || destino.ID == 0 {
		return nil, ErrReservatorioDestinoInvalido
	}
	if origem.ID == destino.ID {
		return nil, ErrReservatorioIgual
	}
	if origem.EmpresaID != empresaID || destino.EmpresaID != empresaID {
		return nil, ErrEmpresaIncompativel
	}
	if origem.Combustivel.ID != destino.Combustivel.ID {
		return nil, ErrCombustivelDiferente
	}
	if quantidade <= 0 {
		return nil, ErrQuantidadeInvalida
	}
	if !origem.Ativo || !destino.Ativo {
		return nil, reservatorio.ErrReservatorioInativo
	}
	if quantidade > origem.NivelAtual {
		return nil, reservatorio.ErrNivelInsuficiente
	}
	if destino.NivelAtual+quantidade > destino.Capacidade {
		return nil, reservatorio.ErrCapacidadeExcedida
	}

	_ = origem.Saida(quantidade)
	_ = destino.Entrada(quantidade)

	return &Abastecimento{
		EmpresaID:             empresaID,
		Tipo:                  TipoTransferencia,
		Data:                  timeutil.Now(),
		Quantidade:            quantidade,
		ReservatorioOrigemID:  origem.ID,
		ReservatorioDestinoID: destino.ID,
	}, nil
}
