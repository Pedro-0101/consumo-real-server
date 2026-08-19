package patrimonio

import (
	"errors"
	"testing"
)

func TestNewPatrimonio(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		p, err := NewPatrimonio(1, "Caminhão", TipoVeiculo, "EXT-1", TipoMedicaoOdometro)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if p.Tipo != TipoVeiculo {
			t.Errorf("tipo esperado VEICULO, obtido %s", p.Tipo)
		}
	})

	t.Run("tipo inválido", func(t *testing.T) {
		_, err := NewPatrimonio(1, "Caminhão", Tipo("NAVE"), "EXT-1", TipoMedicaoOdometro)
		if !errors.Is(err, ErrTipoInvalido) {
			t.Errorf("esperado ErrTipoInvalido, obtido %v", err)
		}
	})

	t.Run("tipo de medição inválido", func(t *testing.T) {
		_, err := NewPatrimonio(1, "Caminhão", TipoVeiculo, "EXT-1", TipoMedicao("PLACA"))
		if !errors.Is(err, ErrTipoMedicaoInvalido) {
			t.Errorf("esperado ErrTipoMedicaoInvalido, obtido %v", err)
		}
	})
}

func TestPatrimonioAtualizar(t *testing.T) {
	p, _ := NewPatrimonio(1, "Caminhão", TipoVeiculo, "EXT-1", TipoMedicaoOdometro)

	if err := p.Atualizar("Trator", TipoMaquina, "EXT-2", TipoMedicaoHorimetro); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if p.Tipo != TipoMaquina || p.TipoMedicao != TipoMedicaoHorimetro {
		t.Errorf("dados não atualizados: tipo=%s medicao=%s", p.Tipo, p.TipoMedicao)
	}

	if err := p.Atualizar("Trator", Tipo("OUTRO"), "EXT-2", TipoMedicaoHorimetro); !errors.Is(err, ErrTipoInvalido) {
		t.Errorf("esperado ErrTipoInvalido, obtido %v", err)
	}
}