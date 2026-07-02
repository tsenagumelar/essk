package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
)

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService, rbacRepo rbac.Repository) {
	auth := middleware.RequireAuth(tokenService)
	group := api.Group("/tenants", auth)

	group.Get("/", middleware.RequirePermission(rbacRepo, "tenants:read"), handler.List)
	group.Post("/", middleware.RequirePermission(rbacRepo, "tenants:create"), handler.Create)
	group.Get("/:id", middleware.RequirePermission(rbacRepo, "tenants:read"), handler.Get)
	group.Patch("/:id", middleware.RequirePermission(rbacRepo, "tenants:update"), handler.Update)
	group.Delete("/:id", middleware.RequirePermission(rbacRepo, "tenants:delete"), handler.Delete)
}
