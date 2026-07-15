package main

import "consumo-real-server/internal/config"

func main() {
	if err := config.Run(); err != nil {
		panic(err)
	}
}
