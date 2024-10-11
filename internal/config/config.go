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
	URI      string `env:"DB_URI" envDefault:"postgres://postgres@localhost:5432/dbname"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD,required"`
}

type Security struct {
	SecretKey string `env:"SECRET_KEY"`
}

type HTTPConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Формирование полного URI с учетом пароля
	cfg.DB.URI = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		"localhost:5432",
		"postgres")

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	return nil
}
