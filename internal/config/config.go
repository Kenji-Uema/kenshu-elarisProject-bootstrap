package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"bootstrap"`
	Version     string `env:"VERSION" envDefault:"dev"`
}

type MongoConfig struct {
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
	Host     string `env:"MONGO_HOST"`
	Database string `env:"MONGO_DATABASE"`
}

type TelemetryConfig struct {
	OTLPEndpoint   string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:31879"`
	OTLPInsecure   bool   `env:"OTEL_EXPORTER_OTLP_INSECURE" envDefault:"true"`
	HealthEndpoint string `env:"OTEL_EXPORTER_OTLP_HEALTH_ENDPOINT" envDefault:"localhost:32019/health"`
}

type PhotosVolumeConfig struct {
	Path string `env:"PHOTOS_VOLUME_PATH" envDefault:"/tmp"`
}

func LoadConfig[C AppConfig | MongoConfig | TelemetryConfig | PhotosVolumeConfig]() C {
	var c C
	if err := env.Parse(&c); err != nil {
		slog.Error("parse env config", "error", err)
	}

	return c
}
