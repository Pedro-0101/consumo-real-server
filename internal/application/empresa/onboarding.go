package empresa

import (
	"context"

	domainempresa "consumo-real-server/internal/domain/empresa"
	domainusuario "consumo-real-server/internal/domain/usuario"
)

// OnboardingRepository agrupa a criação de uma empresa e do seu primeiro
// administrador em uma única transação. O callback de criação do administrador
// é invocado dentro da transação, já com o ID da empresa criada.
type OnboardingRepository interface {
	CriarEmpresaComAdministrador(
		ctx context.Context,
		e *domainempresa.Empresa,
		criarAdministrador func(empresaID int64) (*domainusuario.Usuario, error),
	) error
}
