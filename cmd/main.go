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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	if err := application.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
