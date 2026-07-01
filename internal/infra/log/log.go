package log

import (
	"log/slog"
	"os"

	"github.com/Kenji-Uema/bootstrap/internal/config"
)

func NewLogger(config config.AppConfig) *slog.Logger {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return slog.New(h).With(
		"app", hostname,
		"service.name", config.ServiceName,
		"service.version", config.Version,
		"service.namespace", config.ServiceNamespace,
		"service.instance.id", hostname,
	)
}
