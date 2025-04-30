package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	UserServiceEndpoint string `env:"USER_SERVICE_ENDPOINT,required"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if cerr := env.Parse(&cfg); cerr != nil {
		return nil, cerr
	}

	return &cfg, nil
}
