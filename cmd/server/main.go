// @title API Consumo Real
// @version 1.0
// @description API de controle de consumo real de combustível. Todas as rotas sob /api exigem autenticação via Bearer Token, exceto /api/auth/login.
// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Cole apenas o token JWT (sem o prefixo "Bearer").
package main

import (
	"fmt"
	"os"

	"consumo-real-server/internal/config"
)

func main() {
	if err := config.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "falha ao iniciar o servidor: %v\n", err)
		os.Exit(1)
	}
}
