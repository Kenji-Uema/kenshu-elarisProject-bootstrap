package mq

import (
	"context"
	"log/slog"

	"github.com/Kenji-Uema/bootstrap/internal/config"
)

func BootstrapRabbitmq(ctx context.Context, configs config.Configs) (func(), func(), error) {
	rabbitMqClient, err := NewRabbitMqConnection(ctx, configs.RabbitMqConfig)
	if err != nil {
		return nil, nil, err
	}

	rabbitmqProducer, err := NewRabbitmqProducer(rabbitMqClient)
	if err != nil {
		if closeErr := rabbitMqClient.Close(); closeErr != nil {
			slog.Error("failed to close rabbitmq connection", "error", closeErr)
		}
		return nil, nil, err
	}

	exchanges := []config.ExchangeConfig{
		configs.RabbitMqConfig.Exchanges.Cleaning,
		configs.RabbitMqConfig.Exchanges.TimeEvent,
		configs.RabbitMqConfig.Exchanges.Invoice,
		configs.RabbitMqConfig.Exchanges.Payment,
		configs.RabbitMqConfig.Exchanges.Communication,
	}

	for _, exchange := range exchanges {
		if err := rabbitmqProducer.DeclareExchange(exchange); err != nil {
			if closeErr := rabbitmqProducer.CloseChannel(); closeErr != nil {
				slog.Error("failed to close rabbitmq channel", "error", closeErr)
			}
			if closeErr := rabbitMqClient.Close(); closeErr != nil {
				slog.Error("failed to close rabbitmq connection", "error", closeErr)
			}
			return nil, nil, err
		}
	}

	return func() {
			if err := rabbitmqProducer.CloseChannel(); err != nil {
				slog.Error("failed to close rabbitmq channel", "error", err)
			}
		},
		func() {
			if err := rabbitMqClient.Close(); err != nil {
				slog.Error("failed to close rabbitmq connection", "error", err)
			}
		},
		nil
}
