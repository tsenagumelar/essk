package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
	"github.com/tsenagumelar/essk/services/backend/internal/validator"
)

type Handler struct {
	service   Service
	rbacRepo  rbac.Repository
	validator *validator.Validator
}

func NewHandler(service Service, rbacRepo rbac.Repository, validator *validator.Validator) Handler {
	return Handler{service: service, rbacRepo: rbacRepo, validator: validator}
}

func (h Handler) List(c *fiber.Ctx) error {
	_, scopeTenantID, _, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.List(c.Context(), scopeTenantID)
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
	_, scopeTenantID, _, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.Get(c.Context(), id, scopeTenantID)
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
	claims, _, isSuperAdmin, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if !isSuperAdmin {
		return apperrors.New("FORBIDDEN", fiber.StatusForbidden, "Forbidden")
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
	claims, scopeTenantID, _, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	result, err := h.service.Update(c.Context(), id, req, scopeTenantID, claims.UserID)
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
	claims, _, isSuperAdmin, ok, err := h.scope(c)
	if err != nil {
		return err
	}
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if !isSuperAdmin {
		return apperrors.New("FORBIDDEN", fiber.StatusForbidden, "Forbidden")
	}
	if err := h.service.Delete(c.Context(), id, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"deleted": true}, nil)
}

func (h Handler) scope(c *fiber.Ctx) (authn.Claims, *uuid.UUID, bool, bool, error) {
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return authn.Claims{}, nil, false, false, nil
	}
	isSuperAdmin, err := h.rbacRepo.UserHasRoleCode(c.Context(), claims.UserID, "super_admin")
	if err != nil {
		return authn.Claims{}, nil, false, false, err
	}
	if isSuperAdmin {
		return claims, nil, true, true, nil
	}
	return claims, claims.TenantID, false, true, nil
}
