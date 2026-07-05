package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
)

func New(cfg config.Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if !cfg.Logging.Pretty {
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("service", cfg.Service.Name).
			Str("env", cfg.Service.Env).
			Logger()
	}

	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(writer).
		With().
		Timestamp().
		Str("service", cfg.Service.Name).
		Str("env", cfg.Service.Env).
		Logger()
}
