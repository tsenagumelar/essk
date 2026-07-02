package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
)

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService) {
	group := api.Group("/auth")

	group.Post("/login", handler.Login)
	group.Post("/refresh", handler.Refresh)
	group.Post("/logout", middleware.RequireAuth(tokenService), handler.Logout)
	group.Get("/me", middleware.RequireAuth(tokenService), handler.Me)
}
