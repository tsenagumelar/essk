package rbac

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

func (h Handler) ListPermissions(c *fiber.Ctx) error {
	result, err := h.service.ListPermissions(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) ListRoles(c *fiber.Ctx) error {
	var tenantID *uuid.UUID
	if raw := c.Query("tenant_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid tenant_id", nil)
		}
		tenantID = &parsed
	}

	result, err := h.service.ListRoles(c.Context(), tenantID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) GetRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	result, err := h.service.GetRole(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) CreateRole(c *fiber.Ctx) error {
	req := CreateRoleRequest{}
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

	result, err := h.service.CreateRole(c.Context(), req, claims.UserID)
	if err != nil {
		return err
	}
	return response.Created(c, "Created", result)
}

func (h Handler) UpdateRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	req := UpdateRoleRequest{}
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

	result, err := h.service.UpdateRole(c.Context(), id, req, claims.UserID)
	if err != nil {
		return err
	}
	return response.OK(c, "OK", result, nil)
}

func (h Handler) DeleteRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if err := h.service.DeleteRole(c.Context(), id, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"deleted": true}, nil)
}

func (h Handler) AssignPermission(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	req := AssignPermissionRequest{}
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
	if err := h.service.AssignPermission(c.Context(), roleID, req, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"assigned": true}, nil)
}

func (h Handler) RemovePermission(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	permissionID, err := uuid.Parse(c.Params("permission_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid permission id", nil)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if err := h.service.RemovePermission(c.Context(), roleID, permissionID, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"removed": true}, nil)
}

func (h Handler) AssignRoleToUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}
	req := AssignRoleRequest{}
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
	if err := h.service.AssignRoleToUser(c.Context(), userID, req, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"assigned": true}, nil)
}

func (h Handler) RemoveRoleFromUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}
	roleID, err := uuid.Parse(c.Params("role_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid role id", nil)
	}
	claims, ok := authn.ClaimsFromContext(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if err := h.service.RemoveRoleFromUser(c.Context(), userID, roleID, claims.UserID); err != nil {
		return err
	}
	return response.OK(c, "OK", fiber.Map{"removed": true}, nil)
}
