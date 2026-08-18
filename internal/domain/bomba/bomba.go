package bomba

import (
	"errors"
	"strings"

	"consumo-real-server/internal/shared"
)

var (
	ErrNomeObrigatorio         = errors.New("nome é obrigatório")
	ErrEmpresaObrigatoria      = errors.New("empresa é obrigatória")
	ErrLocalObrigatorio        = errors.New("bomba fixa deve estar vinculada a um local")
	ErrReservatorioObrigatorio = errors.New("reservatório é obrigatório")
	ErrReservatorioInvalido    = errors.New("reservatório inválido")
	ErrBicoObrigatorio         = errors.New("nome do bico é obrigatório")
	ErrBicoInvalido            = errors.New("bico inválido")
)

type Bico struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	BombaID int64  `gorm:"not null;index" json:"bombaID"`
	Nome    string `gorm:"size:100;not null" json:"nome"`
	Ativo   bool   `gorm:"not null;default:true" json:"ativo"`
}

type Bomba struct {
	ID                 int64  `gorm:"primaryKey" json:"id"`
	EmpresaID          int64  `gorm:"not null;index" json:"empresaID"`
	LocalID            int64  `gorm:"index" json:"localID"`
	Movel              bool   `gorm:"not null;default:false" json:"movel"`
	Nome               string `gorm:"size:255;not null" json:"nome"`
	Descricao          string `gorm:"type:text" json:"descricao"`
	ReservatorioID     int64  `gorm:"not null;index" json:"reservatorioID"`
	Bicos              []Bico `gorm:"foreignKey:BombaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bicos"`
	Ativo              bool   `gorm:"not null;default:true" json:"ativo"`
	shared.AuditFields `gorm:"embedded;embeddedPrefix:"`
}

func NewBomba(empresaID, localID, reservatorioID int64, movel bool, nome, descricao string) (*Bomba, error) {
	if empresaID <= 0 {
		return nil, ErrEmpresaObrigatoria
	}
	if strings.TrimSpace(nome) == "" {
		return nil, ErrNomeObrigatorio
	}
	if !movel && localID <= 0 {
		return nil, ErrLocalObrigatorio
	}
	if reservatorioID <= 0 {
		return nil, ErrReservatorioObrigatorio
	}

	return &Bomba{
		EmpresaID:      empresaID,
		LocalID:        localID,
		Movel:          movel,
		Nome:           strings.TrimSpace(nome),
		Descricao:      strings.TrimSpace(descricao),
		ReservatorioID: reservatorioID,
		Bicos:          []Bico{},
		Ativo:          true,
	}, nil
}

func (b *Bomba) VincularReservatorio(reservatorioID int64) error {
	if reservatorioID <= 0 {
		return ErrReservatorioInvalido
	}
	b.ReservatorioID = reservatorioID
	return nil
}

func (b *Bomba) AdicionarBico(nome string) (*Bico, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, ErrBicoObrigatorio
	}
	bico := Bico{
		ID:    b.proximoIDBico(),
		Nome:  strings.TrimSpace(nome),
		Ativo: true,
	}
	b.Bicos = append(b.Bicos, bico)
	return &bico, nil
}

func (b *Bomba) DesativarBico(bicoID int64) error {
	for i := range b.Bicos {
		if b.Bicos[i].ID == bicoID {
			b.Bicos[i].Ativo = false
			return nil
		}
	}
	return ErrBicoInvalido
}

func (b *Bomba) TemBicoAtivo(bicoID int64) bool {
	for _, bico := range b.Bicos {
		if bico.ID == bicoID && bico.Ativo {
			return true
		}
	}
	return false
}

func (b *Bomba) Desativar() {
	b.Ativo = false
}

func (b *Bomba) Atualizar(localID, reservatorioID int64, movel bool, nome, descricao string) error {
	if strings.TrimSpace(nome) == "" {
		return ErrNomeObrigatorio
	}
	if !movel && localID <= 0 {
		return ErrLocalObrigatorio
	}
	if reservatorioID <= 0 {
		return ErrReservatorioObrigatorio
	}

	b.LocalID = localID
	b.Movel = movel
	b.Nome = strings.TrimSpace(nome)
	b.Descricao = strings.TrimSpace(descricao)
	b.ReservatorioID = reservatorioID
	return nil
}

func (b *Bomba) proximoIDBico() int64 {
	var maxID int64
	for _, bico := range b.Bicos {
		if bico.ID > maxID {
			maxID = bico.ID
		}
	}
	return maxID + 1
}
