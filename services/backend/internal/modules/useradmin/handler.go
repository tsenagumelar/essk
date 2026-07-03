package useradmin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
	"github.com/tsenagumelar/essk/services/backend/internal/validator"
)

type Handler struct {
	service   Service
	repo      Repository
	validator *validator.Validator
}

func NewHandler(service Service, repo Repository, validator *validator.Validator) Handler {
	return Handler{service: service, repo: repo, validator: validator}
}

func (h Handler) List(c *fiber.Ctx) error {
	claims, scopeTenantID, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	_ = claims
	result, err := h.service.List(c.Context(), scopeTenantID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Create(c *fiber.Ctx) error {
	claims, scopeTenantID, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	req := CreateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	result, err := h.service.Create(c.Context(), req, scopeTenantID, claims.UserID)
	if err != nil {
		return err
	}
	return response.Created(c, "Created", result)
}

func (h Handler) Update(c *fiber.Ctx) error {
	claims, scopeTenantID, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}
	req := UpdateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if validationErrors := h.validator.Struct(req); validationErrors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation Error", validationErrors)
	}
	result, err := h.service.Update(c.Context(), id, req, scopeTenantID, claims.UserID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) Delete(c *fiber.Ctx) error {
	claims, scopeTenantID, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}
	if err := h.service.Delete(c.Context(), id, scopeTenantID, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"deleted": true}, nil)
}

func (h Handler) scope(c *fiber.Ctx) (authn.Claims, *uuid.UUID, bool, error) {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return authn.Claims{}, nil, false, nil
	}
	isSuperAdmin, err := h.repo.UserHasRoleCode(c.Context(), claims.UserID, "super_admin")
	if err != nil {
		return authn.Claims{}, nil, false, err
	}
	if isSuperAdmin {
		return claims, nil, true, nil
	}
	return claims, claims.TenantID, true, nil
}
