package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"consumo-real-server/internal/shared/auth"
)

type Handlers struct {
	Combustivel  *CombustivelHandler
	Usuario      *UsuarioHandler
	Auth         *AuthHandler
	Empresa      *EmpresaHandler
	UnidadeAdmin *UnidadeAdministrativaHandler
	Local        *LocalHandler
	Patrimonio   *PatrimonioHandler
}

func NewRoutes(handlers Handlers, tokens auth.TokenManager) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")

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

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
