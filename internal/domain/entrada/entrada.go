package entrada

import (
	"errors"
	"strings"
	"time"

	"consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/shared"
	"consumo-real-server/internal/shared/timeutil"
)

var (
	ErrEmpresaObrigatorio    = errors.New("empresa é obrigatória")
	ErrFornecedorObrigatorio = errors.New("fornecedor é obrigatório")
	ErrReservatorioInvalido  = errors.New("reservatório é obrigatório")
	ErrEmpresaIncompativel   = errors.New("reservatório pertence a outra empresa")
	ErrQuantidadeInvalida    = errors.New("quantidade deve ser maior que zero")
)

type Entrada struct {
	ID             int64
	EmpresaID      int64
	FornecedorID   int64
	ReservatorioID int64
	Combustivel    combustivel.Combustivel
	Quantidade     float64
	NotaFiscal     string
	Data           time.Time
	Audit          shared.AuditFields
}

func NewEntrada(empresaID, fornecedorID int64, reservatorio *reservatorio.Reservatorio, quantidade float64, notaFiscal string) (*Entrada, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatorio
	}
	if fornecedorID <= 0 {
		return nil, ErrFornecedorObrigatorio
	}
	if reservatorio == nil || reservatorio.ID == 0 {
		return nil, ErrReservatorioInvalido
	}
	if reservatorio.EmpresaID != empresaID {
		return nil, ErrEmpresaIncompativel
	}
	if quantidade <= 0 {
		return nil, ErrQuantidadeInvalida
	}
	if err := reservatorio.Entrada(quantidade); err != nil {
		return nil, err
	}

	return &Entrada{
		EmpresaID:      empresaID,
		FornecedorID:   fornecedorID,
		ReservatorioID: reservatorio.ID,
		Combustivel:    reservatorio.Combustivel,
		Quantidade:     quantidade,
		NotaFiscal:     strings.TrimSpace(notaFiscal),
		Data:           timeutil.Now(),
	}, nil
}
