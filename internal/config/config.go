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
	RabbitMqConfig
}

type AppConfig struct {
	ServiceName  string `env:"SERVICE_NAME" envDefault:"kenshu-elarisProject-bootstrap"`
	Version      string `env:"VERSION"`
	PhotosVolume struct {
		Path string `env:"PHOTOS_VOLUME_PATH,required"`
	}
}

type MongoConfig struct {
	Username    Secret `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password    Secret `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	Host        string `env:"MONGO_HOST,required"`
	Database    string `env:"MONGO_DATABASE,required"`
	Collections struct {
		CottageCollection string `env:"MONGO_COLLECTION_COTTAGE,required"`
		GuestCollection   string `env:"MONGO_COLLECTION_GUEST,required"`
		BookingCollection string `env:"MONGO_COLLECTION_BOOKING,required"`
		InvoiceCollection string `env:"MONGO_COLLECTION_INVOICE,required"`
		ReceiptCollection string `env:"MONGO_COLLECTION_RECEIPT,required"`
		StockCollection   string `env:"MONGO_COLLECTION_STOCK,required"`
	}
}

type RabbitMqConfig struct {
	Username  Secret `env:"RABBITMQ_USERNAME,required"`
	Password  Secret `env:"RABBITMQ_PASSWORD,required"`
	Host      string `env:"RABBITMQ_HOST,required"`
	Port      int    `env:"RABBITMQ_PORT,required"`
	Exchanges struct {
		Cleaning     ExchangeConfig `envPrefix:"CLEANING_EXCHANGE_"`
		TimeEvent    ExchangeConfig `envPrefix:"TIME_EVENT_EXCHANGE_"`
		Invoice      ExchangeConfig `envPrefix:"INVOICE_EXCHANGE_"`
		Payment      ExchangeConfig `envPrefix:"PAYMENT_EXCHANGE_"`
		Notification ExchangeConfig `envPrefix:"NOTIFICATION_EXCHANGE_"`
	}
}

type TelemetryConfig struct {
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,required"`
	OTLPGrpcPort int    `env:"OTEL_EXPORTER_OTLP_GRPC_PORT,required"`
	OTLPInsecure bool   `env:"OTEL_EXPORTER_OTLP_INSECURE,required"`
}

func LoadConfigs() (Configs, error) {
	var cfg Configs
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	slog.Info("config loaded", "config", cfg)

	return cfg, nil
}
