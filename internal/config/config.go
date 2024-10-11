package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

const (
	EnvDBURI     = "DB_URI"
	EnvSecretKey = "SECRET_KEY"
	EnvHTTPPort  = "HTTP_PORT"
	EnvLogLevel  = "LOG_LEVEL"
)

type Config struct {
	DB       DB
	Security Security
	HTTP     HTTPConfig
}

type DB struct {
	URI string `env:"DB_URI" envDefault:"postgres://postgres:postgres@localhost:5432/dbname"`
}

type Security struct {
	SecretKey string `env:"SECRET_KEY,required"`
}

type HTTPConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Security.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}
	if len(c.Security.SecretKey) < 32 {
		return fmt.Errorf("secret key should be at least 32 characters long")
	}

	return nil
}
