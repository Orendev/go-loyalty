package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	Server struct {
		Addr string `env:"RUN_ADDRESS"`
	}
	Database struct {
		URI string `env:"DATABASE_URI"`
	}
	AccrualSystem struct {
		Addr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	}
	Size int
}

var cfg = Config{}

func NewConfig() (*Config, error) {
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Server.Addr) == 0 {
		flag.StringVar(&cfg.Server.Addr, "a", "", "address to run HTTP server")
	}

	if len(cfg.Database.URI) == 0 {
		flag.StringVar(&cfg.Database.URI, "d", "", "URI to database")
	}

	if len(cfg.AccrualSystem.Addr) == 0 {
		flag.StringVar(&cfg.AccrualSystem.Addr, "r", "", "address of the accrual system")
	}

	cfg.Size = 10

	return &cfg, nil
}
