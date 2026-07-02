package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (a *App) requestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		a.log.Info().
			Str("request_id", c.GetRespHeader(fiber.HeaderXRequestID)).
			Str("correlation_id", c.GetRespHeader(correlationIDHeader)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Int64("latency_ms", latency.Milliseconds()).
			Msg("request completed")

		return err
	}
}
