package rbac

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
)

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService, repo Repository) {
	auth := middleware.RequireAuth(tokenService)

	api.Get("/permissions", auth, middleware.RequirePermission(repo, "permissions:read"), handler.ListPermissions)

	roles := api.Group("/roles", auth)
	roles.Get("/", middleware.RequirePermission(repo, "roles:read"), handler.ListRoles)
	roles.Post("/", middleware.RequirePermission(repo, "roles:create"), handler.CreateRole)
	roles.Get("/:id", middleware.RequirePermission(repo, "roles:read"), handler.GetRole)
	roles.Patch("/:id", middleware.RequirePermission(repo, "roles:update"), handler.UpdateRole)
	roles.Delete("/:id", middleware.RequirePermission(repo, "roles:delete"), handler.DeleteRole)
	roles.Post("/:id/permissions", middleware.RequirePermission(repo, "roles:manage_permissions"), handler.AssignPermission)
	roles.Delete("/:id/permissions/:permission_id", middleware.RequirePermission(repo, "roles:manage_permissions"), handler.RemovePermission)

	users := api.Group("/users", auth)
	users.Post("/:id/roles", middleware.RequirePermission(repo, "users:manage_roles"), handler.AssignRoleToUser)
	users.Delete("/:id/roles/:role_id", middleware.RequirePermission(repo, "users:manage_roles"), handler.RemoveRoleFromUser)
}
