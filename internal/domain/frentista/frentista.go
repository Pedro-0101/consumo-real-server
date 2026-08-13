package frentista

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio    = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria = errors.New("empresa é obrigatória")
)

type Frentista struct {
	ID        int64 `gorm:"primaryKey"`
	EmpresaID int64 `gorm:"not null;index"`
	UsuarioID int64 `gorm:"index"`
	Nome      string `gorm:"size:255;not null"`
	Matricula string `gorm:"size:50;index"`
	Ativo     bool `gorm:"not null;default:true"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewFrentista(empresaID int64, nome, matricula string) (*Frentista, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}

	return &Frentista{
		EmpresaID: empresaID,
		Nome:      strings.TrimSpace(nome),
		Matricula: strings.TrimSpace(matricula),
		Ativo:     true,
	}, nil
}

func (f *Frentista) VincularUsuario(usuarioID int64) {
	f.UsuarioID = usuarioID
}

func (f *Frentista) Desativar() {
	f.Ativo = false
}
