package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/usecase"
)

type Handler struct {
	usecase usecase.Usecase
}

func New(usecase usecase.Usecase) Handler {
	return Handler{usecase: usecase}
}

func (h Handler) Health(c *fiber.Ctx) error {
	return c.JSON(h.usecase.Health(c.Context()))
}

func (h Handler) Upstreams(c *fiber.Ctx) error {
	return c.JSON(h.usecase.Upstreams(c.Context()))
}
