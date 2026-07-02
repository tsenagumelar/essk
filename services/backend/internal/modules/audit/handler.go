package audit

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	var tenantID *string
	if raw := c.Query("tenant_id"); raw != "" {
		tenantID = &raw
	}
	var actorUserID *string
	if raw := c.Query("actor_user_id"); raw != "" {
		actorUserID = &raw
	}

	result, meta, err := h.service.List(c.Context(), ListQuery{
		TenantID:     tenantID,
		ActorUserID:  actorUserID,
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, meta)
}

func (h Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid audit log id", nil)
	}

	result, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}
