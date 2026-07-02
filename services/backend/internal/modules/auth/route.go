package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
)

type RateLimiters struct {
	Login   fiber.Handler
	Refresh fiber.Handler
}

func RegisterRoutes(api fiber.Router, handler Handler, tokenService authn.TokenService, limiters ...RateLimiters) {
	group := api.Group("/auth")

	loginHandlers := []fiber.Handler{handler.Login}
	refreshHandlers := []fiber.Handler{handler.Refresh}
	if len(limiters) > 0 {
		if limiters[0].Login != nil {
			loginHandlers = append([]fiber.Handler{limiters[0].Login}, loginHandlers...)
		}
		if limiters[0].Refresh != nil {
			refreshHandlers = append([]fiber.Handler{limiters[0].Refresh}, refreshHandlers...)
		}
	}

	group.Post("/login", loginHandlers...)
	group.Post("/refresh", refreshHandlers...)
	group.Post("/logout", middleware.RequireAuth(tokenService), handler.Logout)
	group.Get("/me", middleware.RequireAuth(tokenService), handler.Me)
}
