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
	ID                      int64  `gorm:"primaryKey" json:"id"`
	EmpresaID               int64  `gorm:"not null;index" json:"empresaID"`
	UnidadeAdministrativaID int64  `gorm:"not null;index" json:"unidadeAdministrativaID"`
	Nome                    string `gorm:"size:255;not null" json:"nome"`
	Descricao               string `gorm:"type:text" json:"descricao"`
	Ativo                   bool   `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields      `gorm:"embedded;embeddedPrefix:"`
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

func (l *Local) Atualizar(unidadeAdministrativaID int64, nome, descricao string) error {
	if unidadeAdministrativaID <= 0 {
		return ErrUnidadeAdministrativaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return ErrNomeObrigatorio
	}

	l.UnidadeAdministrativaID = unidadeAdministrativaID
	l.Nome = strings.TrimSpace(nome)
	l.Descricao = strings.TrimSpace(descricao)
	return nil
}
