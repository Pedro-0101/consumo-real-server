package frentista

import (
	"errors"
	"testing"
)

func TestNewFrentista(t *testing.T) {
	f, err := NewFrentista(1, "João", "M-001")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if f.UsuarioID != nil {
		t.Errorf("usuarioID deveria ser nil, obtido %v", f.UsuarioID)
	}
}

func TestFrentistaNomeObrigatorio(t *testing.T) {
	_, err := NewFrentista(1, "  ", "M-001")
	if !errors.Is(err, ErrNomeObrigatorio) {
		t.Errorf("esperado ErrNomeObrigatorio, obtido %v", err)
	}
}

func TestVincularUsuario(t *testing.T) {
	f, _ := NewFrentista(1, "João", "M-001")

	f.VincularUsuario(42)
	if f.UsuarioID == nil || *f.UsuarioID != 42 {
		t.Errorf("esperado 42, obtido %v", f.UsuarioID)
	}

	f.VincularUsuario(0)
	if f.UsuarioID != nil {
		t.Errorf("esperado nil após desvincular, obtido %v", f.UsuarioID)
	}
}