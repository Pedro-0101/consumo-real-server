package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"consumo-real-server/internal/shared/auth"
)

type Handlers struct {
	Combustivel *CombustivelHandler
	Usuario     *UsuarioHandler
	Auth        *AuthHandler
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

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
