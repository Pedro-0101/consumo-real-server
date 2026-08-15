package routes

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"consumo-real-server/docs"
	"consumo-real-server/internal/shared/auth"
)

type Handlers struct {
	Combustivel   *CombustivelHandler
	Usuario       *UsuarioHandler
	Auth          *AuthHandler
	Empresa       *EmpresaHandler
	UnidadeAdmin  *UnidadeAdministrativaHandler
	Local         *LocalHandler
	Patrimonio    *PatrimonioHandler
	Reservatorio  *ReservatorioHandler
	Bomba         *BombaHandler
	Frentista     *FrentistaHandler
	Fornecedor    *FornecedorHandler
	Preco         *PrecoHandler
	Ordem         *OrdemAbastecimentoHandler
	Medicao       *MedicaoHandler
	Entrada       *EntradaHandler
	Abastecimento *AbastecimentoHandler
}

func NewRoutes(handlers Handlers, tokens auth.TokenManager) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/swagger", swaggerUIHandler).Methods("GET")
	r.HandleFunc("/swagger/", swaggerUIHandler).Methods("GET")
	r.HandleFunc("/swagger/doc.json", swaggerDocHandler).Methods("GET")

	// Rotas públicas
	r.HandleFunc("/api/auth/login", handlers.Auth.login).Methods("POST")

	// Rotas protegidas
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewareAutenticacao(tokens))
	api.HandleFunc("/auth/me", handlers.Usuario.me).Methods("GET")

	api.HandleFunc("/combustiveis", handlers.Combustivel.list).Methods("GET")
	api.HandleFunc("/combustiveis", handlers.Combustivel.create).Methods("POST")
	api.HandleFunc("/combustiveis/{id}", handlers.Combustivel.get).Methods("GET")
	api.HandleFunc("/combustiveis/{id}", handlers.Combustivel.update).Methods("PUT")
	api.HandleFunc("/combustiveis/{id}", handlers.Combustivel.delete).Methods("DELETE")

	api.HandleFunc("/usuarios", handlers.Usuario.list).Methods("GET")
	api.HandleFunc("/usuarios", handlers.Usuario.create).Methods("POST")
	api.HandleFunc("/usuarios/{id}", handlers.Usuario.get).Methods("GET")
	api.HandleFunc("/usuarios/{id}", handlers.Usuario.update).Methods("PUT")
	api.HandleFunc("/usuarios/{id}", handlers.Usuario.delete).Methods("DELETE")
	api.HandleFunc("/usuarios/{id}/senha", handlers.Usuario.changePassword).Methods("PATCH")

	api.HandleFunc("/empresas", handlers.Empresa.list).Methods("GET")
	api.HandleFunc("/empresas", handlers.Empresa.create).Methods("POST")
	api.HandleFunc("/empresas/{id}", handlers.Empresa.get).Methods("GET")
	api.HandleFunc("/empresas/{id}", handlers.Empresa.update).Methods("PUT")
	api.HandleFunc("/empresas/{id}", handlers.Empresa.delete).Methods("DELETE")

	api.HandleFunc("/unidades-administrativas", handlers.UnidadeAdmin.list).Methods("GET")
	api.HandleFunc("/unidades-administrativas", handlers.UnidadeAdmin.create).Methods("POST")
	api.HandleFunc("/unidades-administrativas/{id}", handlers.UnidadeAdmin.get).Methods("GET")
	api.HandleFunc("/unidades-administrativas/{id}", handlers.UnidadeAdmin.update).Methods("PUT")
	api.HandleFunc("/unidades-administrativas/{id}", handlers.UnidadeAdmin.delete).Methods("DELETE")

	api.HandleFunc("/locais", handlers.Local.list).Methods("GET")
	api.HandleFunc("/locais", handlers.Local.create).Methods("POST")
	api.HandleFunc("/locais/{id}", handlers.Local.get).Methods("GET")
	api.HandleFunc("/locais/{id}", handlers.Local.update).Methods("PUT")
	api.HandleFunc("/locais/{id}", handlers.Local.delete).Methods("DELETE")

	api.HandleFunc("/patrimonios", handlers.Patrimonio.list).Methods("GET")
	api.HandleFunc("/patrimonios", handlers.Patrimonio.create).Methods("POST")
	api.HandleFunc("/patrimonios/{id}", handlers.Patrimonio.get).Methods("GET")
	api.HandleFunc("/patrimonios/{id}", handlers.Patrimonio.update).Methods("PUT")
	api.HandleFunc("/patrimonios/{id}", handlers.Patrimonio.delete).Methods("DELETE")

	api.HandleFunc("/reservatorios", handlers.Reservatorio.list).Methods("GET")
	api.HandleFunc("/reservatorios", handlers.Reservatorio.create).Methods("POST")
	api.HandleFunc("/reservatorios/{id}", handlers.Reservatorio.get).Methods("GET")
	api.HandleFunc("/reservatorios/{id}", handlers.Reservatorio.update).Methods("PUT")
	api.HandleFunc("/reservatorios/{id}", handlers.Reservatorio.delete).Methods("DELETE")

	api.HandleFunc("/bombas", handlers.Bomba.list).Methods("GET")
	api.HandleFunc("/bombas", handlers.Bomba.create).Methods("POST")
	api.HandleFunc("/bombas/{id}", handlers.Bomba.get).Methods("GET")
	api.HandleFunc("/bombas/{id}", handlers.Bomba.update).Methods("PUT")
	api.HandleFunc("/bombas/{id}", handlers.Bomba.delete).Methods("DELETE")
	api.HandleFunc("/bombas/{id}/bicos", handlers.Bomba.adicionarBico).Methods("POST")
	api.HandleFunc("/bombas/{id}/bicos", handlers.Bomba.desativarBico).Methods("DELETE")

	api.HandleFunc("/frentistas", handlers.Frentista.list).Methods("GET")
	api.HandleFunc("/frentistas", handlers.Frentista.create).Methods("POST")
	api.HandleFunc("/frentistas/{id}", handlers.Frentista.get).Methods("GET")
	api.HandleFunc("/frentistas/{id}", handlers.Frentista.update).Methods("PUT")
	api.HandleFunc("/frentistas/{id}", handlers.Frentista.delete).Methods("DELETE")
	api.HandleFunc("/frentistas/{id}/usuario", handlers.Frentista.vincularUsuario).Methods("PATCH")

	api.HandleFunc("/fornecedores", handlers.Fornecedor.list).Methods("GET")
	api.HandleFunc("/fornecedores", handlers.Fornecedor.create).Methods("POST")
	api.HandleFunc("/fornecedores/{id}", handlers.Fornecedor.get).Methods("GET")
	api.HandleFunc("/fornecedores/{id}", handlers.Fornecedor.update).Methods("PUT")
	api.HandleFunc("/fornecedores/{id}", handlers.Fornecedor.delete).Methods("DELETE")

	api.HandleFunc("/precos", handlers.Preco.list).Methods("GET")
	api.HandleFunc("/precos", handlers.Preco.create).Methods("POST")
	api.HandleFunc("/precos/{id}", handlers.Preco.get).Methods("GET")
	api.HandleFunc("/precos/{id}", handlers.Preco.update).Methods("PUT")
	api.HandleFunc("/precos/{id}", handlers.Preco.delete).Methods("DELETE")

	api.HandleFunc("/ordens-abastecimento", handlers.Ordem.list).Methods("GET")
	api.HandleFunc("/ordens-abastecimento", handlers.Ordem.create).Methods("POST")
	api.HandleFunc("/ordens-abastecimento/{id}", handlers.Ordem.get).Methods("GET")
	api.HandleFunc("/ordens-abastecimento/{id}", handlers.Ordem.update).Methods("PUT")
	api.HandleFunc("/ordens-abastecimento/{id}", handlers.Ordem.delete).Methods("DELETE")
	api.HandleFunc("/ordens-abastecimento/{id}/autorizar", handlers.Ordem.autorizar).Methods("PATCH")

	api.HandleFunc("/medicoes", handlers.Medicao.list).Methods("GET")
	api.HandleFunc("/medicoes", handlers.Medicao.create).Methods("POST")
	api.HandleFunc("/medicoes/{id}", handlers.Medicao.get).Methods("GET")
	api.HandleFunc("/medicoes/{id}", handlers.Medicao.update).Methods("PUT")
	api.HandleFunc("/medicoes/{id}", handlers.Medicao.delete).Methods("DELETE")

	api.HandleFunc("/entradas", handlers.Entrada.list).Methods("GET")
	api.HandleFunc("/entradas", handlers.Entrada.create).Methods("POST")
	api.HandleFunc("/entradas/{id}", handlers.Entrada.get).Methods("GET")
	api.HandleFunc("/entradas/{id}", handlers.Entrada.update).Methods("PUT")
	api.HandleFunc("/entradas/{id}", handlers.Entrada.delete).Methods("DELETE")

	api.HandleFunc("/abastecimentos", handlers.Abastecimento.list).Methods("GET")
	api.HandleFunc("/abastecimentos", handlers.Abastecimento.create).Methods("POST")
	api.HandleFunc("/abastecimentos/transferencias", handlers.Abastecimento.createTransferencia).Methods("POST")
	api.HandleFunc("/abastecimentos/{id}", handlers.Abastecimento.get).Methods("GET")
	api.HandleFunc("/abastecimentos/{id}", handlers.Abastecimento.update).Methods("PUT")
	api.HandleFunc("/abastecimentos/{id}", handlers.Abastecimento.delete).Methods("DELETE")

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// swaggerDocHandler serve a especificação Swagger 2.0 gerada automaticamente.
func swaggerDocHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(docs.SwaggerJSON)
}

const swaggerUITemplate = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8" />
  <title>API Consumo Real - Documentação</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
  <script id="swagger-spec" type="application/json">__SWAGGER_SPEC__</script>
  <script>
    window.onload = function () {
      var spec = JSON.parse(document.getElementById("swagger-spec").textContent);
      // Usa a mesma origem da página para evitar erros de CORS:
      // o host do spec passa a ser o host do navegador, garantindo same-origin.
      spec.host = window.location.host;
      spec.schemes = [window.location.protocol.replace(":", "")];
      window.ui = SwaggerUIBundle({
        spec: spec,
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout",
        // cola apenas o token no Authorize: o prefixo "Bearer " é adicionado
        // automaticamente em toda requisição autenticada.
        requestInterceptor: function (req) {
          var auth = req.headers["Authorization"];
          if (auth && auth.indexOf("Bearer ") !== 0) {
            req.headers["Authorization"] = "Bearer " + auth;
          }
          return req;
        }
      });
    };
  </script>
</body>
</html>`

// swaggerUIHandler serve a interface Swagger UI (Swagger UI é carregada via CDN).
// A especificação é embutida na própria página para evitar requisições
// adicionais (evita problemas de CORS, scheme e rede na carga do spec).
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	spec := strings.ReplaceAll(string(docs.SwaggerJSON), "</", "<\\/")
	_, _ = w.Write([]byte(strings.ReplaceAll(swaggerUITemplate, "__SWAGGER_SPEC__", spec)))
}
