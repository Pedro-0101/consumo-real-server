package reservatorio

import (
	"errors"
	"testing"

	"consumo-real-server/internal/domain/combustivel"
)

func novoCombustivel() combustivel.Combustivel {
	return combustivel.Combustivel{ID: 1, EmpresaID: 1}
}

func TestReservatorioAtualizar(t *testing.T) {
	r, err := NewReservatorio(1, "Tanque A", 10000, 9000, 1000, novoCombustivel())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	t.Run("reduzir capacidade abaixo do nível atual", func(t *testing.T) {
		err := r.Atualizar("Tanque A", 8000, 1000, novoCombustivel())
		if !errors.Is(err, ErrCapacidadeExcedida) {
			t.Errorf("esperado ErrCapacidadeExcedida, obtido %v", err)
		}
	})

	t.Run("manter capacidade compatível", func(t *testing.T) {
		if err := r.Atualizar("Tanque A", 10000, 1000, novoCombustivel()); err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
	})
}

func TestReservatorioEntradaSaida(t *testing.T) {
	r, _ := NewReservatorio(1, "Tanque A", 1000, 100, 10, novoCombustivel())

	if err := r.Entrada(100); err != nil {
		t.Fatalf("entrada falhou: %v", err)
	}
	if r.NivelAtual != 200 {
		t.Errorf("nível esperado 200, obtido %v", r.NivelAtual)
	}

	if err := r.Entrada(901); !errors.Is(err, ErrCapacidadeExcedida) {
		t.Errorf("esperado ErrCapacidadeExcedida, obtido %v", err)
	}

	if err := r.Saida(201); !errors.Is(err, ErrNivelInsuficiente) {
		t.Errorf("esperado ErrNivelInsuficiente, obtido %v", err)
	}
}