package main

import (
	"context"
	"github.com/andrepulo/Calendar/internal/application"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Создаем контекст, который будет отменен при получении сигнала завершения (SIGTERM) или прерывания (Ctrl+C)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel() // Отменяем контекст при выходе из функции main

	// Запускаем приложение и проверяем на наличие ошибок
	if err := application.Run(ctx); err != nil {
		log.Fatal(err) // Логируем ошибку и завершаем программу, если произошла ошибка
	}
}
