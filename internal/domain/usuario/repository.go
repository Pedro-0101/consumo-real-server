package usuario

import (
	"context"
	"errors"
)

// ErrNaoEncontrado é retornado pelos repositórios quando o registro não existe.
var ErrNaoEncontrado = errors.New("usuário não encontrado")

// ErrCredenciaisInvalidas é retornado quando o e-mail ou a senha estão incorretos.
var ErrCredenciaisInvalidas = errors.New("credenciais inválidas")

// ListFilter define os filtros disponíveis para a consulta de usuários.
type ListFilter struct {
	EmpresaID int64
	Papel     Papel
	Ativo     *bool
}

// Repository é o contrato de persistência do agregado Usuario.
type Repository interface {
	Create(ctx context.Context, u *Usuario) error
	Update(ctx context.Context, u *Usuario) error
	FindByID(ctx context.Context, id int64) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	List(ctx context.Context, filter ListFilter) ([]Usuario, error)
}

// PasswordHasher é o contrato de hash e verificação de senhas.
type PasswordHasher interface {
	Hash(senha string) (string, error)
	Verificar(hash, senha string) bool
}
