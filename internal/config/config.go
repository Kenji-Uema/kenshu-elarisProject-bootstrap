package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Secret string

func (s Secret) String() string {
	return "REDACTED"
}

type Configs struct {
	AppConfig
	MongoConfig
	TelemetryConfig
	PhotosVolumeConfig
	RabbitMqConfig
	CleaningExchangeConfig
	TimeEventExchangeConfig
	InvoiceExchangeConfig
	PaymentExchangeConfig
	NotificationExchangeConfig
}

type AppConfig struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"kenshu-elarisProject-bootstrap"`
	Version     string `env:"VERSION" envDefault:"latest"`
}

type MongoConfig struct {
	Username Secret `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password Secret `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	Host     string `env:"MONGO_HOST,required"`
	Database string `env:"MONGO_DATABASE,required"`
}

type RabbitMqConfig struct {
	Username Secret `env:"RABBITMQ_USERNAME,required"`
	Password Secret `env:"RABBITMQ_PASSWORD,required"`
	Host     string `env:"RABBITMQ_HOST,required"`
	Port     int    `env:"RABBITMQ_PORT,required"`
}

type TelemetryConfig struct {
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,required"`
	OTLPGrpcPort int    `env:"OTEL_EXPORTER_OTLP_GRPC_PORT,required"`
	OTLPInsecure bool   `env:"OTEL_EXPORTER_OTLP_INSECURE,required"`
}

type PhotosVolumeConfig struct {
	Path string `env:"PHOTOS_VOLUME_PATH,required"`
}

type CleaningExchangeConfig struct {
	Exchange ExchangeConfig `envPrefix:"CLEANING_EXCHANGE_"`
}

type TimeEventExchangeConfig struct {
	Exchange ExchangeConfig `envPrefix:"TIME_EVENT_EXCHANGE_"`
}

type InvoiceExchangeConfig struct {
	Exchange ExchangeConfig `envPrefix:"INVOICE_EXCHANGE_"`
}

type PaymentExchangeConfig struct {
	Exchange ExchangeConfig `envPrefix:"PAYMENT_EXCHANGE_"`
}

type NotificationExchangeConfig struct {
	Exchange ExchangeConfig `envPrefix:"NOTIFICATION_EXCHANGE_"`
}

func LoadConfigs() (Configs, error) {
	var cfg Configs
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	slog.Info("config loaded", "config", cfg)

	return cfg, nil
}
