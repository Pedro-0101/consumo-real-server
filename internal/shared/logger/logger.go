package logger

import (
	"log/slog"
	"os"
)

const ServiceName = "consumo-real-server"

// New cria um logger estruturado em JSON com atributos padrão do serviço.
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler).With("servico", ServiceName)
}
