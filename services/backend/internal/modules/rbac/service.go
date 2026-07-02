package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) ListPermissions(ctx context.Context) ([]PermissionResponse, error) {
	permissions, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		result = append(result, toPermissionResponse(permission))
	}
	return result, nil
}

func (s Service) ListRoles(ctx context.Context, tenantID *uuid.UUID) ([]RoleResponse, error) {
	roles, err := s.repo.ListRoles(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]RoleResponse, 0, len(roles))
	for _, role := range roles {
		result = append(result, toRoleResponse(role))
	}
	return result, nil
}

func (s Service) GetRole(ctx context.Context, id uuid.UUID) (RoleResponse, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return RoleResponse{}, mapNotFound(err)
	}
	return toRoleResponse(role), nil
}

func (s Service) CreateRole(ctx context.Context, req CreateRoleRequest, actorID uuid.UUID) (RoleResponse, error) {
	var tenantID *uuid.UUID
	if req.TenantID != nil && *req.TenantID != "" {
		parsed, err := uuid.Parse(*req.TenantID)
		if err != nil {
			return RoleResponse{}, apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid tenant_id")
		}
		tenantID = &parsed
	}

	role := Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsSystem:    req.IsSystem,
		IsActive:    true,
	}

	if err := s.repo.CreateRole(ctx, role, actorID, s.now().UTC()); err != nil {
		return RoleResponse{}, err
	}
	return toRoleResponse(role), nil
}

func (s Service) UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleRequest, actorID uuid.UUID) (RoleResponse, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return RoleResponse{}, mapNotFound(err)
	}

	role.Name = req.Name
	role.Description = req.Description
	role.IsActive = req.IsActive

	if err := s.repo.UpdateRole(ctx, role, actorID, s.now().UTC()); err != nil {
		return RoleResponse{}, mapNotFound(err)
	}
	return toRoleResponse(role), nil
}

func (s Service) DeleteRole(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	return mapNotFound(s.repo.DeleteRole(ctx, id, actorID, s.now().UTC()))
}

func (s Service) AssignPermission(ctx context.Context, roleID uuid.UUID, req AssignPermissionRequest, actorID uuid.UUID) error {
	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		return apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid permission_id")
	}
	return s.repo.AssignPermission(ctx, roleID, permissionID, actorID, s.now().UTC())
}

func (s Service) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID, actorID uuid.UUID) error {
	return s.repo.RemovePermission(ctx, roleID, permissionID, actorID, s.now().UTC())
}

func (s Service) AssignRoleToUser(ctx context.Context, userID uuid.UUID, req AssignRoleRequest, actorID uuid.UUID) error {
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid role_id")
	}
	return s.repo.AssignRoleToUser(ctx, userID, roleID, actorID, s.now().UTC())
}

func (s Service) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, actorID uuid.UUID) error {
	return s.repo.RemoveRoleFromUser(ctx, userID, roleID, actorID, s.now().UTC())
}

func mapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return apperrors.New("NOT_FOUND", fiber.StatusNotFound, "Resource not found")
	}
	return err
}

func toPermissionResponse(permission Permission) PermissionResponse {
	return PermissionResponse{
		ID:          permission.ID.String(),
		Code:        permission.Code,
		Name:        permission.Name,
		Description: permission.Description,
		IsActive:    permission.IsActive,
	}
}

func toRoleResponse(role Role) RoleResponse {
	var tenantID *string
	if role.TenantID != nil {
		value := role.TenantID.String()
		tenantID = &value
	}
	return RoleResponse{
		ID:          role.ID.String(),
		TenantID:    tenantID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
	}
}
