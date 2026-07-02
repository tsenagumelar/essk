package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
)

func RequireAuth(tokens authn.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		rawToken, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(rawToken) == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		claims, err := tokens.ParseAccessToken(rawToken)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		c.Locals(authn.ClaimsLocalKey, claims)
		return c.Next()
	}
}
