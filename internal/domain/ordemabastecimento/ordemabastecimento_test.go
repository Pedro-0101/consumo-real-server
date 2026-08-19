package ordemabastecimento

import (
	"errors"
	"testing"
)

func TestNewOrdemAbastecimento(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		o, err := NewOrdemAbastecimento(1, 10, "OA-001", 100, nil)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if o.Status != StatusAberta {
			t.Errorf("status esperado ABERTA, obtido %s", o.Status)
		}
	})

	t.Run("número obrigatório", func(t *testing.T) {
		_, err := NewOrdemAbastecimento(1, 10, "  ", 100, nil)
		if !errors.Is(err, ErrNumeroObrigatorio) {
			t.Errorf("esperado ErrNumeroObrigatorio, obtido %v", err)
		}
	})

	t.Run("quantidade inválida", func(t *testing.T) {
		_, err := NewOrdemAbastecimento(1, 10, "OA-002", 0, nil)
		if !errors.Is(err, ErrQuantidadeInvalida) {
			t.Errorf("esperado ErrQuantidadeInvalida, obtido %v", err)
		}
	})
}

func TestCicloDeVida(t *testing.T) {
	o, _ := NewOrdemAbastecimento(1, 10, "OA-001", 100, nil)

	if err := o.Autorizar(); err != nil {
		t.Fatalf("autorizar falhou: %v", err)
	}
	if o.Status != StatusAutorizada {
		t.Errorf("status esperado AUTORIZADA, obtido %s", o.Status)
	}

	if err := o.RegistrarAbastecimento(40); err != nil {
		t.Fatalf("registrar abastecimento falhou: %v", err)
	}
	if o.QuantidadeAbastecida != 40 || o.Status != StatusAutorizada {
		t.Errorf("estado inesperado: quantidade=%v status=%s", o.QuantidadeAbastecida, o.Status)
	}

	if err := o.RegistrarAbastecimento(60); err != nil {
		t.Fatalf("registrar abastecimento final falhou: %v", err)
	}
	if o.Status != StatusConcluida {
		t.Errorf("status esperado CONCLUIDA, obtido %s", o.Status)
	}

	if err := o.RegistrarAbastecimento(1); !errors.Is(err, ErrStatusInvalido) {
		t.Errorf("esperado ErrStatusInvalido, obtido %v", err)
	}
	if err := o.Cancelar(); !errors.Is(err, ErrStatusInvalido) {
		t.Errorf("esperado ErrStatusInvalido ao cancelar concluída, obtido %v", err)
	}
}