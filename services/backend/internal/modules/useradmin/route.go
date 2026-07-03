package useradmin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
)

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService, rbacRepo rbac.Repository) {
	auth := middleware.RequireAuth(tokenService)
	group := api.Group("/users", auth)

	group.Get("/", middleware.RequirePermission(rbacRepo, "users:read"), handler.List)
	group.Post("/", middleware.RequirePermission(rbacRepo, "users:create"), handler.Create)
	group.Patch("/:id", middleware.RequirePermission(rbacRepo, "users:update"), handler.Update)
	group.Delete("/:id", middleware.RequirePermission(rbacRepo, "users:delete"), handler.Delete)
}
