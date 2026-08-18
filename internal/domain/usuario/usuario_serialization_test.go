package usuario

import (
	"encoding/json"
	"testing"
	"time"

	"consumo-real-server/internal/shared"
)

func TestUsuarioJSONKeysCamelCase(t *testing.T) {
	now := time.Now()
	u := &Usuario{
		ID:        1,
		EmpresaID: 10,
		Nome:      "Admin",
		Email:     "a@b.com",
		SenhaHash: "secret",
		Papel:     PapelAdminBase,
		Ativo:     true,
		AuditFields: shared.AuditFields{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: 1,
			UpdatedBy: 1,
		},
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal falhou: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal falhou: %v", err)
	}

	for _, chave := range []string{"id", "empresaID", "nome", "email", "papel", "ativo", "createdAt", "updatedAt", "createdBy", "updatedBy"} {
		if _, ok := m[chave]; !ok {
			t.Errorf("chave %q ausente no JSON: %s", chave, string(data))
		}
	}

	if _, ok := m["SenhaHash"]; ok {
		t.Errorf("chave PascalCase inesperada no JSON: %s", string(data))
	}
	if _, ok := m["senhaHash"]; ok {
		t.Errorf("campo senhaHash não deveria ser serializado: %s", string(data))
	}
}
