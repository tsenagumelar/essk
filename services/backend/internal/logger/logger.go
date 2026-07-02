package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

func New(cfg config.Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	if !cfg.Logging.Pretty {
		return zerolog.New(os.Stdout).With().Timestamp().Str("app", cfg.App.Name).Str("env", cfg.App.Env).Logger()
	}

	return zerolog.New(writer).With().Timestamp().Str("app", cfg.App.Name).Str("env", cfg.App.Env).Logger()
}
