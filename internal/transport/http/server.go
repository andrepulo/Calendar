package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrepulo/Calendar/internal/config"
)

// HealthCheckHandler - простой обработчик для проверки работы сервера
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// StartServer запускает HTTP-сервер в отдельной горутине
func StartServer(ctx context.Context, cfg *config.HTTPConfig, handler http.Handler) (*http.Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	server := &http.Server{Handler: handler}

	go func() {
		err := server.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return server, nil
}

// GracefulShutdown корректно завершает работу сервера
func GracefulShutdown(ctx context.Context, server *http.Server) error {
	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("принудительное завершение: %w", err)
	}

	fmt.Println("Сервер корректно завершил работу")
	return nil
}

// WaitForShutdown ожидает сигнала завершения и запускает корректное завершение сервера
func WaitForShutdown(ctx context.Context, server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		fmt.Println("Контекст отменен, завершение работы сервера...")
	case <-quit:
		fmt.Println("Получен сигнал завершения, завершение работы сервера...")
	}

	if err := GracefulShutdown(ctx, server); err != nil {
		fmt.Printf("Ошибка при завершении работы сервера: %v\n", err)
	}
}
