package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
	"github.com/tsenagumelar/essk/services/backend/internal/validator"
)

type Handler struct {
	service   Service
	validator *validator.Validator
}

func NewHandler(service Service, validator *validator.Validator) Handler {
	return Handler{service: service, validator: validator}
}

func (h Handler) List(c *fiber.Ctx) error {
	result, err := h.service.List(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid tenant id", nil)
	}
	result, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Create(c *fiber.Ctx) error {
	req := CreateTenantRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.Create(c.Context(), req, claims.UserID)
	if err != nil {
		return err
	}
	return response.Created(c, "Created", result)
}

func (h Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid tenant id", nil)
	}
	req := UpdateTenantRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.Update(c.Context(), id, req, claims.UserID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid tenant id", nil)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if err := h.service.Delete(c.Context(), id, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"deleted": true}, nil)
}
