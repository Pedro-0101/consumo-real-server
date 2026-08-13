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
	ID                    int64
	EmpresaID             int64
	LocalID               int64
	BombaID               int64
	BicoID                int64
	Tipo                  Tipo
	Data                  time.Time
	Quantidade            float64
	PrecoUnitario         float64
	ValorTotal            float64
	Odometro              float64
	Horimetro             float64
	Combustivel           combustivel.Combustivel
	PatrimonioID          int64
	FrentistaID           int64
	ReservatorioOrigemID  int64
	ReservatorioDestinoID int64
	Audit                 shared.AuditFields
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
