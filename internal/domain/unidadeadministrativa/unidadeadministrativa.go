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
	ID                      int64
	EmpresaID               int64
	UnidadeAdministrativaID int64
	Nome                    string
	Tipo                    Tipo
	Ativo                   bool
	Audit                   shared.AuditFields
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
