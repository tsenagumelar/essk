package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
)

type Service struct {
	repo  Repository
	audit *audit.Service
	now   func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) WithAudit(auditService audit.Service) Service {
	s.audit = &auditService
	return s
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

func (s Service) GetRole(ctx context.Context, id uuid.UUID, scopeTenantID *uuid.UUID) (RoleResponse, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return RoleResponse{}, mapNotFound(err)
	}
	if err := requireRoleScope(role, scopeTenantID); err != nil {
		return RoleResponse{}, err
	}
	return toRoleResponse(role), nil
}

func (s Service) CreateRole(ctx context.Context, req CreateRoleRequest, scopeTenantID *uuid.UUID, actorID uuid.UUID) (RoleResponse, error) {
	var tenantID *uuid.UUID
	if req.TenantID != nil && *req.TenantID != "" {
		parsed, err := uuid.Parse(*req.TenantID)
		if err != nil {
			return RoleResponse{}, apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid tenant_id")
		}
		tenantID = &parsed
	}
	if scopeTenantID != nil {
		if tenantID != nil && *tenantID != *scopeTenantID {
			return RoleResponse{}, apperrors.New("FORBIDDEN", fiber.StatusForbidden, "Forbidden")
		}
		tenantID = scopeTenantID
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
	_ = s.writeAudit(ctx, actorID, "role.create", "role", role.ID.String(), map[string]any{"code": role.Code})
	return toRoleResponse(role), nil
}

func (s Service) UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleRequest, scopeTenantID *uuid.UUID, actorID uuid.UUID) (RoleResponse, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return RoleResponse{}, mapNotFound(err)
	}
	if err := requireRoleScope(role, scopeTenantID); err != nil {
		return RoleResponse{}, err
	}

	role.Name = req.Name
	role.Description = req.Description
	role.IsActive = req.IsActive

	if err := s.repo.UpdateRole(ctx, role, actorID, s.now().UTC()); err != nil {
		return RoleResponse{}, mapNotFound(err)
	}
	_ = s.writeAudit(ctx, actorID, "role.update", "role", role.ID.String(), map[string]any{"code": role.Code})
	return toRoleResponse(role), nil
}

func (s Service) DeleteRole(ctx context.Context, id uuid.UUID, scopeTenantID *uuid.UUID, actorID uuid.UUID) error {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return mapNotFound(err)
	}
	if err := requireRoleScope(role, scopeTenantID); err != nil {
		return err
	}
	if err := mapNotFound(s.repo.DeleteRole(ctx, id, actorID, s.now().UTC())); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "role.delete", "role", id.String(), nil)
	return nil
}

func (s Service) AssignPermission(ctx context.Context, roleID uuid.UUID, req AssignPermissionRequest, actorID uuid.UUID) error {
	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		return apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid permission_id")
	}
	if err := s.repo.AssignPermission(ctx, roleID, permissionID, actorID, s.now().UTC()); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "role.permission.assign", "role", roleID.String(), map[string]any{"permission_id": permissionID.String()})
	return nil
}

func (s Service) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID, actorID uuid.UUID) error {
	if err := s.repo.RemovePermission(ctx, roleID, permissionID, actorID, s.now().UTC()); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "role.permission.remove", "role", roleID.String(), map[string]any{"permission_id": permissionID.String()})
	return nil
}

func (s Service) AssignRoleToUser(ctx context.Context, userID uuid.UUID, req AssignRoleRequest, actorID uuid.UUID) error {
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid role_id")
	}
	if err := s.repo.AssignRoleToUser(ctx, userID, roleID, actorID, s.now().UTC()); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "user.role.assign", "user", userID.String(), map[string]any{"role_id": roleID.String()})
	return nil
}

func (s Service) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, actorID uuid.UUID) error {
	if err := s.repo.RemoveRoleFromUser(ctx, userID, roleID, actorID, s.now().UTC()); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "user.role.remove", "user", userID.String(), map[string]any{"role_id": roleID.String()})
	return nil
}

func (s Service) writeAudit(ctx context.Context, actorID uuid.UUID, action string, resourceType string, resourceID string, metadata map[string]any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(ctx, audit.Event{
		ActorUserID:  &actorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   &resourceID,
		Metadata:     metadata,
	})
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

func requireRoleScope(role Role, scopeTenantID *uuid.UUID) error {
	if scopeTenantID == nil {
		return nil
	}
	if role.TenantID == nil || *role.TenantID != *scopeTenantID {
		return apperrors.New("FORBIDDEN", fiber.StatusForbidden, "Forbidden")
	}
	return nil
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
