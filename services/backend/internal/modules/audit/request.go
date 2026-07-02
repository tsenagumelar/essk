package audit

import "github.com/gofiber/fiber/v2"

func RequestMetadata(c *fiber.Ctx) (*string, *string) {
	ip := c.IP()
	userAgent := c.Get(fiber.HeaderUserAgent)
	return &ip, &userAgent
}
