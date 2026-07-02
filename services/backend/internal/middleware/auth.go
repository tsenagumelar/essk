package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
)

type PermissionChecker interface {
	UserHasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error)
}

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

func RequirePermission(repo PermissionChecker, permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := authn.ClaimsFromContext(c)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		allowed, err := repo.UserHasPermission(c.Context(), claims.UserID, permissionCode)
		if err != nil {
			return err
		}
		if !allowed {
			return response.Error(c, fiber.StatusForbidden, "Forbidden", []map[string]string{
				{"code": "FORBIDDEN", "permission": permissionCode},
			})
		}

		return c.Next()
	}
}
