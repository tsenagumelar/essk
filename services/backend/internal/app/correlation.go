package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const correlationIDHeader = "X-Correlation-ID"

func correlationID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(correlationIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Locals("correlation_id", id)
		c.Set(correlationIDHeader, id)
		return c.Next()
	}
}
