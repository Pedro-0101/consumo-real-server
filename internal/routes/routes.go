package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRoutes() (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	return r, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
