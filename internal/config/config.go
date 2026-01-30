package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	ServiceName string `env:"SERVICE_NAME" required:"true"`
	Version     string `env:"VERSION" required:"true"`
}

type MongoConfig struct {
	Username string `env:"MONGO_INITDB_ROOT_USERNAME" required:"true"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD" required:"true"`
	Host     string `env:"MONGO_HOST" required:"true"`
	Database string `env:"MONGO_DATABASE" required:"true"`
}

type TelemetryConfig struct {
	OTLPEndpoint   string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" required:"true"`
	OTLPInsecure   bool   `env:"OTEL_EXPORTER_OTLP_INSECURE" required:"true"`
	HealthEndpoint string `env:"OTEL_EXPORTER_OTLP_HEALTH_ENDPOINT" required:"true"`
}

type PhotosVolumeConfig struct {
	Path string `env:"PHOTOS_VOLUME_PATH" required:"true"`
}

func LoadConfig[C AppConfig | MongoConfig | TelemetryConfig | PhotosVolumeConfig]() (C, error) {
	var c C
	if err := env.Parse(&c); err != nil {
		slog.Error("parse env config", "error", err)
		return c, err
	}

	return c, nil
}
