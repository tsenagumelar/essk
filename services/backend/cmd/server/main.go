package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsenagumelar/essk/services/backend/internal/app"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
	"github.com/tsenagumelar/essk/services/backend/internal/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg)

	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create application")
	}

	go func() {
		if err := application.Listen(); err != nil {
			log.Fatal().Err(err).Msg("server stopped unexpectedly")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown cleanly")
	}
}
