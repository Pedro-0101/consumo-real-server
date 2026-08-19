package entrada

import (
	"errors"
	"testing"

	"consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/domain/reservatorio"
)

func novoReservatorio(t *testing.T) *reservatorio.Reservatorio {
	t.Helper()
	c := combustivel.Combustivel{ID: 1, EmpresaID: 1}
	r, err := reservatorio.NewReservatorio(1, "Tanque A", 10000, 1000, 100, c)
	if err != nil {
		t.Fatalf("falha ao criar reservatório: %v", err)
	}
	r.ID = 1
	return r
}

func TestNewEntrada(t *testing.T) {
	t.Run("sucesso com nota fiscal", func(t *testing.T) {
		e, err := NewEntrada(1, 5, novoReservatorio(t), 500, "NF 123")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if e.NotaFiscal == nil || *e.NotaFiscal != "NF 123" {
			t.Errorf("nota fiscal inesperada: %v", e.NotaFiscal)
		}
	})

	t.Run("nota fiscal vazia vira nil", func(t *testing.T) {
		e, err := NewEntrada(1, 5, novoReservatorio(t), 500, "  ")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if e.NotaFiscal != nil {
			t.Errorf("nota fiscal deveria ser nil, obtido %v", e.NotaFiscal)
		}
	})

	t.Run("quantidade inválida", func(t *testing.T) {
		_, err := NewEntrada(1, 5, novoReservatorio(t), 0, "")
		if !errors.Is(err, ErrQuantidadeInvalida) {
			t.Errorf("esperado ErrQuantidadeInvalida, obtido %v", err)
		}
	})
}

func TestEntradaAtualizarNotaFiscal(t *testing.T) {
	e, _ := NewEntrada(1, 5, novoReservatorio(t), 500, "NF 123")

	e.AtualizarNotaFiscal("NF 456")
	if e.NotaFiscal == nil || *e.NotaFiscal != "NF 456" {
		t.Errorf("nota fiscal inesperada: %v", e.NotaFiscal)
	}

	e.AtualizarNotaFiscal("")
	if e.NotaFiscal != nil {
		t.Errorf("nota fiscal deveria ser nil, obtido %v", e.NotaFiscal)
	}
}