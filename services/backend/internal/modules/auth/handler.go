package auth

import (
	"github.com/gofiber/fiber/v2"
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

func (h Handler) Login(c *fiber.Ctx) error {
	req := LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}

	result, err := h.service.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return response.OK(c, "OK", result, nil)
}

func (h Handler) Refresh(c *fiber.Ctx) error {
	req := RefreshRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}

	result, err := h.service.Refresh(c.Context(), req)
	if err != nil {
		return err
	}

	return response.OK(c, "OK", result, nil)
}

func (h Handler) Logout(c *fiber.Ctx) error {
	req := LogoutRequest{}
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

	if err := h.service.Logout(c.Context(), req, claims.UserID); err != nil {
		return err
	}

	return response.OK(c, "OK", fiber.Map{"logged_out": true}, nil)
}

func (h Handler) Me(c *fiber.Ctx) error {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	result, err := h.service.Me(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return response.OK(c, "OK", result, nil)
}
