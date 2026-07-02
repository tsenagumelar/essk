package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
)

func (a *App) health(c *fiber.Ctx) error {
	return response.OK(c, "OK", fiber.Map{
		"app":     a.cfg.App.Name,
		"env":     a.cfg.App.Env,
		"version": a.cfg.App.Version,
		"status":  "healthy",
	}, nil)
}

func (a *App) ready(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	checks := fiber.Map{
		"database": "not_configured",
		"redis":    "not_configured",
	}

	ready := true

	if a.db != nil {
		if err := a.db.Ping(ctx); err != nil {
			checks["database"] = "unhealthy"
			ready = false
		} else {
			checks["database"] = "healthy"
		}
	}

	if a.redis != nil {
		if err := a.redis.Ping(ctx).Err(); err != nil {
			checks["redis"] = "unhealthy"
			ready = false
		} else {
			checks["redis"] = "healthy"
		}
	}

	if !ready {
		return response.Error(c, fiber.StatusServiceUnavailable, "Service Unavailable", checks)
	}

	return response.OK(c, "OK", fiber.Map{
		"status": "ready",
		"checks": checks,
	}, nil)
}

func (a *App) version(c *fiber.Ctx) error {
	return response.OK(c, "OK", fiber.Map{
		"version": a.cfg.App.Version,
	}, nil)
}
