package audit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
)

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService, permissionChecker middleware.PermissionChecker) {
	auth := middleware.RequireAuth(tokenService)
	group := api.Group("/audit-logs", auth, middleware.RequirePermission(permissionChecker, "audit_logs:read"))

	group.Get("/", handler.List)
	group.Get("/:id", handler.Get)
}
