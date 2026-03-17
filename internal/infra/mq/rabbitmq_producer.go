package mq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kenji-Uema/bootstrap/internal/config"
	"github.com/Kenji-Uema/bootstrap/internal/port"
)

type rabbitmqProducer struct {
	*RabbitMqChannel
	exchangeName string
	exchangeKind string
}

func NewRabbitmqProducer(rabbitmqConnection *RabbitMqConnection) (port.MqProducer, error) {
	paymentProducer := rabbitmqProducer{
		RabbitMqChannel: NewRabbitMqChannel(rabbitmqConnection),
	}

	if err := paymentProducer.openChannel(); err != nil {
		return nil, err
	}

	return &paymentProducer, nil
}

func (p *rabbitmqProducer) DeclareExchange(config config.ExchangeConfig) error {
	p.exchangeName = config.Name
	p.exchangeKind = config.Kind
	if config.Kind == "" {
		slog.Warn("exchange kind not specified, defaulting to 'direct'")
		p.exchangeKind = "direct"
	}

	if p.channel == nil || p.channel.IsClosed() {
		if err := p.reopenChannel(context.Background()); err != nil {
			return err
		}
	}

	if err := p.channel.ExchangeDeclare(p.exchangeName, p.exchangeKind,
		config.Durable, config.AutoDelete, config.Internal,
		config.NoWait, nil); err != nil {

		return fmt.Errorf("declare exchange %q: %w", config.Name, err)
	}

	return nil
}
