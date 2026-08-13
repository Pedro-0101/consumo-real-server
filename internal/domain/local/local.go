package local

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio                  = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria               = errors.New("empresa é obrigatória")
	ErrUnidadeAdministrativaObrigatoria = errors.New("unidade administrativa é obrigatória")
)

type Local struct {
	ID                      int64 `gorm:"primaryKey"`
	EmpresaID               int64 `gorm:"not null;index"`
	UnidadeAdministrativaID int64 `gorm:"not null;index"`
	Nome                    string `gorm:"size:255;not null"`
	Descricao               string `gorm:"type:text"`
	Ativo                   bool `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewLocal(empresaID, unidadeAdministrativaID int64, nome, descricao string) (*Local, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if unidadeAdministrativaID <= 0 {
		return nil, ErrUnidadeAdministrativaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}

	return &Local{
		EmpresaID:               empresaID,
		UnidadeAdministrativaID: unidadeAdministrativaID,
		Nome:                    strings.TrimSpace(nome),
		Descricao:               strings.TrimSpace(descricao),
		Ativo:                   true,
	}, nil
}

func (l *Local) Desativar() {
	l.Ativo = false
}
