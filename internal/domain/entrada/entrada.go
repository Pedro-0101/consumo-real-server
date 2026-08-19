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
	ID                 int64                   `gorm:"primaryKey" json:"id"`
	EmpresaID          int64                   `gorm:"not null;uniqueIndex:uk_empresa_notafiscal,priority:1" json:"empresaID"`
	FornecedorID       int64                   `gorm:"not null;index" json:"fornecedorID"`
	ReservatorioID     int64                   `gorm:"not null;index" json:"reservatorioID"`
	CombustivelID      int64                   `gorm:"not null;index" json:"combustivelID"`
	Combustivel        combustivel.Combustivel `gorm:"foreignKey:CombustivelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"combustivel"`
	Quantidade         float64                 `gorm:"not null" json:"quantidade"`
	NotaFiscal         *string                 `gorm:"size:60;uniqueIndex:uk_empresa_notafiscal,priority:2" json:"notaFiscal"`
	Data               time.Time               `gorm:"not null;index" json:"data"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
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
		CombustivelID:  reservatorio.Combustivel.ID,
		Combustivel:    reservatorio.Combustivel,
		Quantidade:     quantidade,
		NotaFiscal:     normalizarNotaFiscal(notaFiscal),
		Data:           timeutil.Now(),
	}, nil
}

// AtualizarNotaFiscal atualiza a nota fiscal da entrada, normalizando o valor.
// Um valor vazio desvincula a nota fiscal (NULL no banco).
func (e *Entrada) AtualizarNotaFiscal(notaFiscal string) {
	e.NotaFiscal = normalizarNotaFiscal(notaFiscal)
}

// normalizarNotaFiscal converte o valor para *string, retornando nil quando vazio.
// Assim o índice único (empresa_id, nota_fiscal) aceita múltiplos NULLs.
func normalizarNotaFiscal(notaFiscal string) *string {
	nf := strings.TrimSpace(notaFiscal)
	if nf == "" {
		return nil
	}
	return &nf
}
