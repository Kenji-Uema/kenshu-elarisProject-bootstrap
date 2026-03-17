package port

import (
	"github.com/Kenji-Uema/bootstrap/internal/config"
)

type MqProducer interface {
	DeclareExchange(config config.ExchangeConfig) error
	CloseChannel() error
}
