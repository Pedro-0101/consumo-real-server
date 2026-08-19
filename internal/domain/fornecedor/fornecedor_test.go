package fornecedor

import (
	"errors"
	"testing"
)

func TestNewFornecedor(t *testing.T) {
	t.Run("sucesso sem CNPJ", func(t *testing.T) {
		f, err := NewFornecedor(1, "Posto Aurora", "")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if f.CNPJ != nil {
			t.Errorf("CNPJ deveria ser nil, obtido %v", f.CNPJ)
		}
	})

	t.Run("CNPJ válido é normalizado para dígitos", func(t *testing.T) {
		f, err := NewFornecedor(1, "Posto Aurora", "11.222.333/0001-81")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if f.CNPJ == nil || *f.CNPJ != "11222333000181" {
			t.Errorf("CNPJ esperado 11222333000181, obtido %v", f.CNPJ)
		}
	})

	t.Run("CNPJ inválido", func(t *testing.T) {
		_, err := NewFornecedor(1, "Posto Aurora", "00000000000000")
		if !errors.Is(err, ErrCNPJInvalido) {
			t.Errorf("esperado ErrCNPJInvalido, obtido %v", err)
		}
	})

	t.Run("nome obrigatório", func(t *testing.T) {
		_, err := NewFornecedor(1, "  ", "")
		if !errors.Is(err, ErrNomeObrigatorio) {
			t.Errorf("esperado ErrNomeObrigatorio, obtido %v", err)
		}
	})
}

func TestFornecedorAtualizar(t *testing.T) {
	f, _ := NewFornecedor(1, "Posto Aurora", "")

	t.Run("atualiza CNPJ válido", func(t *testing.T) {
		if err := f.Atualizar("Posto Aurora 2", "11.222.333/0001-81"); err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if f.CNPJ == nil || *f.CNPJ != "11222333000181" {
			t.Errorf("CNPJ esperado 11222333000181, obtido %v", f.CNPJ)
		}
	})

	t.Run("CNPJ inválido", func(t *testing.T) {
		if err := f.Atualizar("Posto Aurora 2", "00000000000000"); !errors.Is(err, ErrCNPJInvalido) {
			t.Errorf("esperado ErrCNPJInvalido, obtido %v", err)
		}
	})

	t.Run("CNPJ vazio desvincula", func(t *testing.T) {
		if err := f.Atualizar("Posto Aurora 2", ""); err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if f.CNPJ != nil {
			t.Errorf("CNPJ deveria ser nil, obtido %v", f.CNPJ)
		}
	})
}