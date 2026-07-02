package product

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
	claims, ok := authn.ClaimsFromContext(c)
	if !ok || claims.TenantID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.List(c.Context(), *claims.TenantID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Get(c *fiber.Ctx) error {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok || claims.TenantID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product id", nil)
	}
	result, err := h.service.Get(c.Context(), *claims.TenantID, id)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Create(c *fiber.Ctx) error {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok || claims.TenantID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	req := CreateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	result, err := h.service.Create(c.Context(), *claims.TenantID, req, claims.UserID)
	if err != nil {
		return err
	}
	return response.Created(c, "Created", result)
}

func (h Handler) Update(c *fiber.Ctx) error {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok || claims.TenantID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product id", nil)
	}
	req := UpdateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	result, err := h.service.Update(c.Context(), *claims.TenantID, id, req, claims.UserID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Delete(c *fiber.Ctx) error {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok || claims.TenantID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product id", nil)
	}
	if err := h.service.Delete(c.Context(), *claims.TenantID, id, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"deleted": true}, nil)
}
