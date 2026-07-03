package useradmin

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
	authmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/auth"
)

type Service struct {
	repo   Repository
	audit  *audit.Service
	hasher authmodule.PasswordHasher
	now    func() time.Time
}

func NewService(repo Repository, hasher authmodule.PasswordHasher) Service {
	return Service{repo: repo, hasher: hasher, now: time.Now}
}

func (s Service) WithAudit(auditService audit.Service) Service {
	s.audit = &auditService
	return s
}

func (s Service) List(ctx context.Context, scopeTenantID *uuid.UUID) ([]UserResponse, error) {
	users, err := s.repo.List(ctx, scopeTenantID)
	if err != nil {
		return nil, err
	}
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response, err := s.toResponse(ctx, user)
		if err != nil {
			return nil, err
		}
		result = append(result, response)
	}
	return result, nil
}

func (s Service) Create(ctx context.Context, req CreateUserRequest, scopeTenantID *uuid.UUID, actorID uuid.UUID) (UserResponse, error) {
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return UserResponse{}, apperrors.New("VALIDATION_ERROR", fiber.StatusBadRequest, "Invalid tenant_id")
	}
	if scopeTenantID != nil && *scopeTenantID != tenantID {
		return UserResponse{}, apperrors.New("FORBIDDEN", fiber.StatusForbidden, "Forbidden")
	}
	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return UserResponse{}, err
	}
	user := User{
		ID:       uuid.New(),
		TenantID: &tenantID,
		Email:    req.Email,
		Name:     req.Name,
		Status:   "active",
		IsActive: true,
	}
	now := s.now().UTC()
	if err := s.repo.Create(ctx, user, passwordHash, actorID, now); err != nil {
		return UserResponse{}, err
	}
	if err := s.repo.ReplaceRoles(ctx, user.ID, parseRoleIDs(req.RoleIDs), actorID, now); err != nil {
		return UserResponse{}, err
	}
	_ = s.writeAudit(ctx, actorID, "user.create", user.ID.String(), map[string]any{"email": user.Email})
	return s.toResponse(ctx, user)
}

func (s Service) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest, scopeTenantID *uuid.UUID, actorID uuid.UUID) (UserResponse, error) {
	user, err := s.repo.Get(ctx, id, scopeTenantID)
	if err != nil {
		return UserResponse{}, mapNotFound(err)
	}
	user.Name = req.Name
	user.Status = req.Status
	user.IsActive = req.IsActive
	now := s.now().UTC()
	if err := s.repo.Update(ctx, user, actorID, now); err != nil {
		return UserResponse{}, mapNotFound(err)
	}
	if err := s.repo.ReplaceRoles(ctx, user.ID, parseRoleIDs(req.RoleIDs), actorID, now); err != nil {
		return UserResponse{}, err
	}
	_ = s.writeAudit(ctx, actorID, "user.update", user.ID.String(), map[string]any{"status": user.Status})
	return s.toResponse(ctx, user)
}

func (s Service) Delete(ctx context.Context, id uuid.UUID, scopeTenantID *uuid.UUID, actorID uuid.UUID) error {
	if _, err := s.repo.Get(ctx, id, scopeTenantID); err != nil {
		return mapNotFound(err)
	}
	if err := mapNotFound(s.repo.Delete(ctx, id, actorID, s.now().UTC())); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "user.delete", id.String(), nil)
	return nil
}

func (s Service) toResponse(ctx context.Context, user User) (UserResponse, error) {
	roleIDs, err := s.repo.ListRoleIDs(ctx, user.ID)
	if err != nil {
		return UserResponse{}, err
	}
	var tenantID *string
	if user.TenantID != nil {
		value := user.TenantID.String()
		tenantID = &value
	}
	result := UserResponse{
		ID:       user.ID.String(),
		TenantID: tenantID,
		Email:    user.Email,
		Name:     user.Name,
		Status:   user.Status,
		IsActive: user.IsActive,
		RoleIDs:  make([]string, 0, len(roleIDs)),
	}
	for _, roleID := range roleIDs {
		result.RoleIDs = append(result.RoleIDs, roleID.String())
	}
	return result, nil
}

func (s Service) writeAudit(ctx context.Context, actorID uuid.UUID, action string, resourceID string, metadata map[string]any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(ctx, audit.Event{
		ActorUserID:  &actorID,
		Action:       action,
		ResourceType: "user",
		ResourceID:   &resourceID,
		Metadata:     metadata,
	})
}

func parseRoleIDs(raw []string) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(raw))
	for _, value := range raw {
		parsed, err := uuid.Parse(value)
		if err == nil {
			result = append(result, parsed)
		}
	}
	return result
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
