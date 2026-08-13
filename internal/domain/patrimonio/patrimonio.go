package patrimonio

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

const (
	TipoVeiculo     = "VEICULO"
	TipoFerramenta  = "FERRAMENTA"
	TipoMaquina     = "MAQUINA"
	TipoEquipamento = "EQUIPAMENTO"
)

type TipoMedicao string

const (
	TipoMedicaoOdometro  TipoMedicao = "ODOMETRO"
	TipoMedicaoHorimetro TipoMedicao = "HORIMETRO"
)

var (
	ErrNomeObrigatorio     = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria  = errors.New("empresa é obrigatória")
	ErrTipoObrigatorio     = errors.New("tipo é obrigatório")
	ErrTipoMedicaoInvalido = errors.New("tipo de medição inválido")
)

type Patrimonio struct {
	ID                      int64 `gorm:"primaryKey"`
	EmpresaID               int64 `gorm:"not null;index"`
	UnidadeAdministrativaID int64 `gorm:"index"`
	Nome                    string `gorm:"size:255;not null"`
	Descricao               string `gorm:"type:text"`
	Tipo                    string `gorm:"type:varchar(30);not null;index"`
	TipoMedicao             TipoMedicao `gorm:"type:varchar(20);not null"`
	CodigoExterno           string `gorm:"size:100;index"`
	Atributos               map[string]string `gorm:"type:jsonb;serializer:json"`
	Ativo                   bool `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewPatrimonio(empresaID int64, nome, tipo, codigoExterno string, tipoMedicao TipoMedicao) (*Patrimonio, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if strings.TrimSpace(tipo) == "" {
		return nil, ErrTipoObrigatorio
	}
	if !tipoMedicao.isValid() {
		return nil, ErrTipoMedicaoInvalido
	}

	return &Patrimonio{
		EmpresaID:     empresaID,
		Nome:          strings.TrimSpace(nome),
		Tipo:          strings.TrimSpace(tipo),
		TipoMedicao:   tipoMedicao,
		CodigoExterno: strings.TrimSpace(codigoExterno),
		Atributos:     make(map[string]string),
		Ativo:         true,
	}, nil
}

func (p *Patrimonio) Alocar(unidadeAdministrativaID int64) {
	p.UnidadeAdministrativaID = unidadeAdministrativaID
}

func (p *Patrimonio) SetAtributo(chave, valor string) {
	if p.Atributos == nil {
		p.Atributos = make(map[string]string)
	}
	p.Atributos[strings.TrimSpace(chave)] = strings.TrimSpace(valor)
}

func (p *Patrimonio) Desativar() {
	p.Ativo = false
}

func (t TipoMedicao) isValid() bool {
	switch t {
	case TipoMedicaoOdometro, TipoMedicaoHorimetro:
		return true
	}
	return false
}
