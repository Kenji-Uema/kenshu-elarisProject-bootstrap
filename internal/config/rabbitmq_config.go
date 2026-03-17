package config

type ExchangeConfig struct {
	Name       string `env:"NAME,required"`
	Kind       string `env:"KIND,required"`
	Durable    bool   `env:"DURABLE" envDefault:"true"`
	AutoDelete bool   `env:"AUTO_DELETE" envDefault:"false"`
	Internal   bool   `env:"INTERNAL" envDefault:"false"`
	NoWait     bool   `env:"NO_WAIT" envDefault:"false"`
}
