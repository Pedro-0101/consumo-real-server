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
