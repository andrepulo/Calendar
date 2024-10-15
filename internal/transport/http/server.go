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
	fmt.Fprintf(w, "Hello, World!") // Отправляем ответ "Hello, World!" клиенту
}

// StartServer запускает HTTP-сервер в отдельной горутине
func StartServer(ctx context.Context, cfg *config.HTTPConfig, handler http.Handler) (*http.Server, error) {
	// Создаем слушатель на указанном порту
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("net: %w", err) // Возвращаем ошибку, если не удалось создать слушатель
	}

	server := &http.Server{Handler: handler} // Создаем новый HTTP-сервер с указанным обработчиком

	go func() {
		// Запускаем сервер и обрабатываем ошибки
		err := server.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err) // Паника, если сервер завершился неожиданно
		}
	}()

	return server, nil // Возвращаем сервер и nil, если все прошло успешно
}

// GracefulShutdown корректно завершает работу сервера
func GracefulShutdown(ctx context.Context, server *http.Server) error {
	// Создаем контекст с таймаутом для завершения работы сервера
	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Завершаем работу сервера
	if err := server.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("принудительное завершение: %w", err) // Возвращаем ошибку, если не удалось корректно завершить работу сервера
	}

	fmt.Println("Сервер корректно завершил работу")
	return nil // Возвращаем nil, если завершение прошло успешно
}

// WaitForShutdown ожидает сигнала завершения и запускает корректное завершение сервера
func WaitForShutdown(ctx context.Context, server *http.Server) {
	quit := make(chan os.Signal, 1)
	// Ожидаем сигналы завершения (SIGINT, SIGTERM)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		fmt.Println("Контекст отменен, завершение работы сервера...") // Сообщение, если контекст был отменен
	case <-quit:
		fmt.Println("Получен сигнал завершения, завершение работы сервера...") // Сообщение, если получен сигнал завершения
	}

	// Корректно завершаем работу сервера
	if err := GracefulShutdown(ctx, server); err != nil {
		fmt.Printf("Ошибка при завершени�� работы сервера: %v\n", err) // Выводим ошибку, если не удалось корректно завершить работу сервера
	}
}
