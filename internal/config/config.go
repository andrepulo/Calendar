package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"log"
)

// Config представляет основную структуру конфигурации приложения.
type Config struct {
	DB       DB         // Конфигурация базы данных
	Security Security   // Конфигурация безопасности
	HTTP     HTTPConfig // Конфигурация HTTP
}

// DB представляет структуру конфигурации базы данных.
type DB struct {
	URI string `env:"DB_URI" envDefault:"postgresql://postgres:password@localhost:5555/auth"` // URI для подключения к базе данных
}

// Security представляет структуру конфигурации безопасности.
type Security struct {
	SecretKey string `env:"SECRET_KEY"` // Секретный ключ для безопасности
}

// HTTPConfig представляет структуру конфигурации HTTP.
type HTTPConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"` // Порт для HTTP-сервера
}

// Parse загружает конфигурацию из переменных окружения.
func Parse() (*Config, error) {
	var cfg Config

	// Загрузка конфигурации из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфигурацию из окружения: %w", err)
	}

	// Логирование загруженной конфигурации
	log.Printf("Конфигурация загружена: %+v\n", cfg)

	// Валидация конфигурации
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("некорректная конфигурация: %w", err)
	}

	return &cfg, nil
}

// validateConfig проверяет загруженную конфигурацию.
func validateConfig(cfg *Config) error {
	if cfg.DB.URI == "" {
		return fmt.Errorf("требуется URI базы данных")
	}
	if cfg.Security.SecretKey == "" {
		return fmt.Errorf("требуется секретный ключ безопасности")
	}
	if cfg.HTTP.Port <= 0 {
		return fmt.Errorf("порт HTTP должен быть положительным числом")
	}
	return nil
}
