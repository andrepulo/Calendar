package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrepulo/Calendar/internal/config"
	"github.com/andrepulo/Calendar/internal/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	loggerConfig := logger.Config{
		Level: os.Getenv(config.EnvLogLevel), // Используем переменную окружения для уровня логирования
	}
	l, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer l.Sync()

	// Простой обработчик для проверки работы сервера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Создание и запуск HTTP-сервера
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
	}

	go func() {
		l.Infof("Starting server on port %d", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatalf("Server forced to shutdown: %v", err)
	}

	l.Info("Server exiting")
}
