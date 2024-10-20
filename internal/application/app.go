package application

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andrepulo/Calendar/internal/config"
	handlers "github.com/andrepulo/Calendar/internal/controllers/http"
	"github.com/andrepulo/Calendar/internal/databases"
	"github.com/andrepulo/Calendar/internal/logger"
	httpTransport "github.com/andrepulo/Calendar/internal/transport/http"
)

func Run(ctx context.Context) error {
	// Загрузка конфигурации
	cfg, err := config.Parse()
	if err != nil {
		return fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	}

	// Инициализация логгера
	loggerConfig := logger.Config{
		Level:     os.Getenv("LOG_LEVEL"),
		LogToFile: os.Getenv("LOG_TO_FILE") == "true",
		FilePath:  os.Getenv("LOG_FILE_PATH"),
	}
	l, err := logger.New(loggerConfig)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}
	defer l.Sync()

	// Инициализация базы данных
	db, closeDB, err := databases.NewDB(&cfg.DB)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer closeDB()

	// Использование нового обработчика для проверки работы сервера
	http.HandleFunc("/", httpTransport.HealthCheckHandler)
	usersDomain := buildUsersDomain(cfg, db)
	authDomain := buildAuthDomain(cfg, usersDomain)
	hs := handlers.NewHandlers(usersDomain.users, authDomain.auth, l)

	// Запуск сервера
	server, err := httpTransport.StartServer(ctx, &cfg.HTTP, http.DefaultServeMux)
	if err != nil {
		return fmt.Errorf("не удалось запустить сервер: %w", err)
	}

	// Ожидание сигнала для graceful shutdown
	httpTransport.WaitForShutdown(ctx, server)

	return nil
}
