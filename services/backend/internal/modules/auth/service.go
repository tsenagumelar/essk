package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
)

type Service struct {
	cfg      config.Config
	repo     Repository
	rbacRepo *rbac.Repository
	hasher   PasswordHasher
	tokenSvc authn.TokenService
	now      func() time.Time
}

func NewService(cfg config.Config, repo Repository, hasher PasswordHasher, tokenSvc authn.TokenService) Service {
	return Service{
		cfg:      cfg,
		repo:     repo,
		hasher:   hasher,
		tokenSvc: tokenSvc,
		now:      time.Now,
	}
}

func (s Service) WithRBAC(repo rbac.Repository) Service {
	s.rbacRepo = &repo
	return s
}

func (s Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return AuthResponse{}, invalidCredentials()
	}

	if !user.IsActive || user.IsDeleted || user.Status != "active" {
		return AuthResponse{}, invalidCredentials()
	}

	ok, err := s.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil || !ok {
		return AuthResponse{}, invalidCredentials()
	}

	now := s.now().UTC()
	if err := s.repo.UpdateLastLogin(ctx, user.ID, now); err != nil {
		return AuthResponse{}, err
	}

	return s.createAuthResponse(ctx, user, now)
}

func (s Service) Refresh(ctx context.Context, req RefreshRequest) (AuthResponse, error) {
	now := s.now().UTC()
	current, err := s.repo.FindRefreshTokenByHash(ctx, HashRefreshToken(req.RefreshToken))
	if err != nil {
		return AuthResponse{}, apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Unauthorized")
	}

	if current.RevokedAt != nil || !current.ExpiresAt.After(now) {
		return AuthResponse{}, apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := s.repo.FindUserByID(ctx, current.UserID)
	if err != nil {
		return AuthResponse{}, apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Unauthorized")
	}

	response, err := s.createAuthResponse(ctx, user, now)
	if err != nil {
		return AuthResponse{}, err
	}

	newToken, err := s.repo.FindRefreshTokenByHash(ctx, HashRefreshToken(response.RefreshToken))
	if err != nil {
		return AuthResponse{}, err
	}
	if err := s.repo.RevokeRefreshToken(ctx, current.ID, user.ID, &newToken.ID, now); err != nil {
		return AuthResponse{}, err
	}

	return response, nil
}

func (s Service) Logout(ctx context.Context, req LogoutRequest, actorID uuid.UUID) error {
	token, err := s.repo.FindRefreshTokenByHash(ctx, HashRefreshToken(req.RefreshToken))
	if err != nil {
		return nil
	}

	now := s.now().UTC()
	return s.repo.RevokeRefreshToken(ctx, token.ID, actorID, nil, now)
}

func (s Service) Me(ctx context.Context, userID uuid.UUID) (UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return UserResponse{}, apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Unauthorized")
	}
	if !user.IsActive || user.IsDeleted || user.Status != "active" {
		return UserResponse{}, apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Unauthorized")
	}
	return toUserResponse(user), nil
}

func (s Service) SeedAdmin(ctx context.Context) error {
	now := s.now().UTC()
	tenantID := uuid.New()
	if err := s.repo.EnsureTenant(ctx, tenantID, s.cfg.Seed.TenantName, s.cfg.Seed.TenantSlug, now); err != nil {
		return err
	}

	existingTenantID, err := s.repo.FindTenantBySlug(ctx, s.cfg.Seed.TenantSlug)
	if err != nil {
		return err
	}

	existingUser, err := s.repo.FindUserByEmail(ctx, s.cfg.Seed.AdminEmail)
	if err == nil {
		if s.rbacRepo == nil {
			return nil
		}
		return s.seedAdminRBAC(ctx, existingUser.ID, existingUser.TenantID, now)
	}
	if !errors.Is(err, ErrNotFound) {
		return err
	}

	passwordHash, err := s.hasher.Hash(s.cfg.Seed.AdminPassword)
	if err != nil {
		return err
	}

	userID := uuid.New()
	user := User{
		ID:           userID,
		TenantID:     &existingTenantID,
		Email:        s.cfg.Seed.AdminEmail,
		Name:         s.cfg.Seed.AdminName,
		PasswordHash: passwordHash,
		Status:       "active",
		IsActive:     true,
	}

	if err := s.repo.EnsureUser(ctx, user, userID, now); err != nil {
		return err
	}

	if s.rbacRepo == nil {
		return nil
	}
	return s.seedAdminRBAC(ctx, userID, &existingTenantID, now)
}

func (s Service) seedAdminRBAC(ctx context.Context, adminUserID uuid.UUID, tenantID *uuid.UUID, now time.Time) error {
	permissions := []rbac.Permission{
		{ID: uuid.New(), Code: "permissions:read", Name: "Read permissions"},
		{ID: uuid.New(), Code: "roles:read", Name: "Read roles"},
		{ID: uuid.New(), Code: "roles:create", Name: "Create roles"},
		{ID: uuid.New(), Code: "roles:update", Name: "Update roles"},
		{ID: uuid.New(), Code: "roles:delete", Name: "Delete roles"},
		{ID: uuid.New(), Code: "roles:manage_permissions", Name: "Manage role permissions"},
		{ID: uuid.New(), Code: "users:manage_roles", Name: "Manage user roles"},
		{ID: uuid.New(), Code: "tenants:read", Name: "Read tenants"},
		{ID: uuid.New(), Code: "tenants:create", Name: "Create tenants"},
		{ID: uuid.New(), Code: "tenants:update", Name: "Update tenants"},
		{ID: uuid.New(), Code: "tenants:delete", Name: "Delete tenants"},
	}

	for _, permission := range permissions {
		if err := s.rbacRepo.EnsurePermission(ctx, permission, &adminUserID, now); err != nil {
			return err
		}
	}

	description := "System administrator role"
	adminRole := rbac.Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        "Administrator",
		Code:        "admin",
		Description: &description,
		IsSystem:    true,
		IsActive:    true,
	}
	if err := s.rbacRepo.EnsureRole(ctx, adminRole, adminUserID, now); err != nil {
		return err
	}

	storedRole, err := s.rbacRepo.FindRoleByCode(ctx, tenantID, "admin")
	if err != nil {
		return err
	}
	if err := s.rbacRepo.AssignRoleToUser(ctx, adminUserID, storedRole.ID, adminUserID, now); err != nil {
		return err
	}

	for _, permission := range permissions {
		storedPermission, err := s.rbacRepo.FindPermissionByCode(ctx, permission.Code)
		if err != nil {
			return err
		}
		if err := s.rbacRepo.AssignPermission(ctx, storedRole.ID, storedPermission.ID, adminUserID, now); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) createAuthResponse(ctx context.Context, user User, now time.Time) (AuthResponse, error) {
	accessToken, accessExpiresAt, err := s.tokenSvc.CreateAccessToken(user.ID, user.TenantID, user.Email, now)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, refreshTokenHash, err := NewRefreshToken()
	if err != nil {
		return AuthResponse{}, err
	}

	refreshExpiresAt := now.Add(s.cfg.Auth.RefreshTokenTTL)
	if err := s.repo.CreateRefreshToken(ctx, RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: refreshExpiresAt,
	}, user.ID, now); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
		User:                  toUserResponse(user),
	}, nil
}

func invalidCredentials() error {
	return apperrors.New("UNAUTHORIZED", fiber.StatusUnauthorized, "Invalid email or password")
}

func toUserResponse(user User) UserResponse {
	var tenantID *string
	if user.TenantID != nil {
		value := user.TenantID.String()
		tenantID = &value
	}

	return UserResponse{
		ID:       user.ID.String(),
		TenantID: tenantID,
		Email:    user.Email,
		Name:     user.Name,
		Status:   user.Status,
	}
}
