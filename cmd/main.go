package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-monitor/internal/app"
	"url-monitor/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	app, err := app.New(ctx, ":"+cfg.AppPort, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("server started on :%s", cfg.AppPort)

		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<- stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Close(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
