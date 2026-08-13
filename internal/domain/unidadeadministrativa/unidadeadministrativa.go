package unidadeadministrativa

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

type Tipo string

const (
	TipoMatriz   Tipo = "MATRIZ"
	TipoFilial   Tipo = "FILIAL"
	TipoDeposito Tipo = "DEPOSITO"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
	ErrTipoInvalido       = errors.New("tipo de unidade administrativa inválido")
)

type UnidadeAdministrativa struct {
	ID                      int64 `gorm:"primaryKey"`
	EmpresaID               int64 `gorm:"not null;index"`
	UnidadeAdministrativaID int64 `gorm:"index"`
	Nome                    string `gorm:"size:255;not null"`
	Tipo                    Tipo `gorm:"type:varchar(20);not null"`
	Ativo                   bool `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewUnidadeAdministrativa(empresaID, unidadeAdministrativaID int64, nome string, tipo Tipo) (*UnidadeAdministrativa, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if !tipo.isValid() {
		return nil, ErrTipoInvalido
	}

	return &UnidadeAdministrativa{
		EmpresaID:               empresaID,
		UnidadeAdministrativaID: unidadeAdministrativaID,
		Nome:                    strings.TrimSpace(nome),
		Tipo:                    tipo,
		Ativo:                   true,
	}, nil
}

func (u *UnidadeAdministrativa) Desativar() {
	u.Ativo = false
}

func (t Tipo) isValid() bool {
	switch t {
	case TipoMatriz, TipoFilial, TipoDeposito:
		return true
	}
	return false
}
