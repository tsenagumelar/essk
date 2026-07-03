package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
)

type Service struct {
	cfg      config.Config
	repo     Repository
	rbacRepo *rbac.Repository
	audit    *audit.Service
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

func (s Service) WithAudit(auditService audit.Service) Service {
	s.audit = &auditService
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
	_ = s.writeAudit(ctx, user, "auth.login", "user", user.ID.String(), nil)

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
	_ = s.writeAudit(ctx, user, "auth.refresh", "refresh_token", current.ID.String(), nil)

	return response, nil
}

func (s Service) Logout(ctx context.Context, req LogoutRequest, actorID uuid.UUID) error {
	token, err := s.repo.FindRefreshTokenByHash(ctx, HashRefreshToken(req.RefreshToken))
	if err != nil {
		return nil
	}

	now := s.now().UTC()
	if err := s.repo.RevokeRefreshToken(ctx, token.ID, actorID, nil, now); err != nil {
		return err
	}
	user, err := s.repo.FindUserByID(ctx, actorID)
	if err == nil {
		_ = s.writeAudit(ctx, user, "auth.logout", "refresh_token", token.ID.String(), nil)
	}
	return nil
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
		return fmt.Errorf("ensure seed tenant %s: %w", s.cfg.Seed.TenantSlug, err)
	}

	existingTenantID, err := s.repo.FindTenantBySlug(ctx, s.cfg.Seed.TenantSlug)
	if err != nil {
		return fmt.Errorf("find seed tenant %s: %w", s.cfg.Seed.TenantSlug, err)
	}

	existingUser, err := s.repo.FindUserByEmail(ctx, s.cfg.Seed.AdminEmail)
	if err == nil {
		if s.rbacRepo == nil {
			return nil
		}
		return s.seedAdminRBAC(ctx, existingUser.ID, existingUser.TenantID, now)
	}
	if !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("find seed admin user %s: %w", s.cfg.Seed.AdminEmail, err)
	}

	passwordHash, err := s.hasher.Hash(s.cfg.Seed.AdminPassword)
	if err != nil {
		return fmt.Errorf("hash seed admin password: %w", err)
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
		return fmt.Errorf("ensure seed admin user %s: %w", s.cfg.Seed.AdminEmail, err)
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
		{ID: uuid.New(), Code: "users:read", Name: "Read users"},
		{ID: uuid.New(), Code: "users:create", Name: "Create users"},
		{ID: uuid.New(), Code: "users:update", Name: "Update users"},
		{ID: uuid.New(), Code: "users:delete", Name: "Delete users"},
		{ID: uuid.New(), Code: "users:manage_roles", Name: "Manage user roles"},
		{ID: uuid.New(), Code: "tenants:read", Name: "Read tenants"},
		{ID: uuid.New(), Code: "tenants:create", Name: "Create tenants"},
		{ID: uuid.New(), Code: "tenants:update", Name: "Update tenants"},
		{ID: uuid.New(), Code: "tenants:delete", Name: "Delete tenants"},
		{ID: uuid.New(), Code: "products:read", Name: "Read products"},
		{ID: uuid.New(), Code: "products:create", Name: "Create products"},
		{ID: uuid.New(), Code: "products:update", Name: "Update products"},
		{ID: uuid.New(), Code: "products:delete", Name: "Delete products"},
		{ID: uuid.New(), Code: "audit_logs:read", Name: "Read audit logs"},
	}

	for _, permission := range permissions {
		if err := s.rbacRepo.EnsurePermission(ctx, permission, &adminUserID, now); err != nil {
			return fmt.Errorf("ensure permission %s: %w", permission.Code, err)
		}
	}

	if err := s.ensureTenantRoles(ctx, tenantID, adminUserID, now, permissions); err != nil {
		return fmt.Errorf("ensure system tenant roles: %w", err)
	}

	superAdminDescription := "Global super administrator role"
	superAdminRole := rbac.Role{
		ID:          uuid.New(),
		TenantID:    nil,
		Name:        "Super Administrator",
		Code:        "super_admin",
		Description: &superAdminDescription,
		IsSystem:    true,
		IsActive:    true,
	}
	if err := s.rbacRepo.EnsureRole(ctx, superAdminRole, adminUserID, now); err != nil {
		return fmt.Errorf("ensure super_admin role: %w", err)
	}

	storedRole, err := s.rbacRepo.FindRoleByCode(ctx, nil, "super_admin")
	if err != nil {
		return fmt.Errorf("find super_admin role: %w", err)
	}
	if err := s.rbacRepo.ReplaceUserRoles(ctx, adminUserID, []uuid.UUID{storedRole.ID}, adminUserID, now); err != nil {
		return fmt.Errorf("replace default admin roles with super_admin: %w", err)
	}

	for _, permission := range permissions {
		storedPermission, err := s.rbacRepo.FindPermissionByCode(ctx, permission.Code)
		if err != nil {
			return fmt.Errorf("find permission %s: %w", permission.Code, err)
		}
		if err := s.rbacRepo.AssignPermission(ctx, storedRole.ID, storedPermission.ID, adminUserID, now); err != nil {
			return fmt.Errorf("assign permission %s to super_admin: %w", permission.Code, err)
		}
	}

	if err := s.seedSampleTenants(ctx, adminUserID, now, permissions); err != nil {
		return fmt.Errorf("seed sample tenants: %w", err)
	}

	return nil
}

func (s Service) ensureTenantRoles(ctx context.Context, tenantID *uuid.UUID, actorID uuid.UUID, now time.Time, permissions []rbac.Permission) error {
	adminDescription := "Tenant administrator role"
	tenantAdminRole := rbac.Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        "Tenant Administrator",
		Code:        "admin",
		Description: &adminDescription,
		IsSystem:    true,
		IsActive:    true,
	}
	if err := s.rbacRepo.EnsureRole(ctx, tenantAdminRole, actorID, now); err != nil {
		return err
	}

	userDescription := "Tenant user role"
	tenantUserRole := rbac.Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        "User",
		Code:        "user",
		Description: &userDescription,
		IsSystem:    true,
		IsActive:    true,
	}
	if err := s.rbacRepo.EnsureRole(ctx, tenantUserRole, actorID, now); err != nil {
		return err
	}

	adminRole, err := s.rbacRepo.FindRoleByCode(ctx, tenantID, "admin")
	if err != nil {
		return err
	}
	userRole, err := s.rbacRepo.FindRoleByCode(ctx, tenantID, "user")
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		storedPermission, err := s.rbacRepo.FindPermissionByCode(ctx, permission.Code)
		if err != nil {
			return err
		}
		if err := s.rbacRepo.AssignPermission(ctx, adminRole.ID, storedPermission.ID, actorID, now); err != nil {
			return err
		}
		if permission.Code == "products:read" {
			if err := s.rbacRepo.AssignPermission(ctx, userRole.ID, storedPermission.ID, actorID, now); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s Service) seedSampleTenants(ctx context.Context, actorID uuid.UUID, now time.Time, permissions []rbac.Permission) error {
	samples := []struct {
		name string
		slug string
	}{
		{name: "Acme Corporation", slug: "acme"},
		{name: "Globex Indonesia", slug: "globex"},
	}

	for _, sample := range samples {
		if err := s.repo.EnsureTenant(ctx, uuid.New(), sample.name, sample.slug, now); err != nil {
			return err
		}
		tenantID, err := s.repo.FindTenantBySlug(ctx, sample.slug)
		if err != nil {
			return err
		}
		if err := s.ensureTenantRoles(ctx, &tenantID, actorID, now, permissions); err != nil {
			return err
		}

		adminRole, err := s.rbacRepo.FindRoleByCode(ctx, &tenantID, "admin")
		if err != nil {
			return err
		}
		userRole, err := s.rbacRepo.FindRoleByCode(ctx, &tenantID, "user")
		if err != nil {
			return err
		}

		if err := s.seedSampleUser(ctx, tenantID, sample.slug+"-admin@essk.local", sample.name+" Admin", "Admin123!", adminRole.ID, actorID, now); err != nil {
			return err
		}
		if err := s.seedSampleUser(ctx, tenantID, sample.slug+"-user@essk.local", sample.name+" User", "User12345!", userRole.ID, actorID, now); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) seedSampleUser(ctx context.Context, tenantID uuid.UUID, email string, name string, password string, roleID uuid.UUID, actorID uuid.UUID, now time.Time) error {
	existing, err := s.repo.FindUserByEmailAndTenant(ctx, email, tenantID)
	if err == nil {
		return s.rbacRepo.AssignRoleToUser(ctx, existing.ID, roleID, actorID, now)
	}
	if !errors.Is(err, ErrNotFound) {
		return err
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}
	user := User{
		ID:           uuid.New(),
		TenantID:     &tenantID,
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		Status:       "active",
		IsActive:     true,
	}
	if err := s.repo.EnsureUser(ctx, user, actorID, now); err != nil {
		return err
	}
	return s.rbacRepo.AssignRoleToUser(ctx, user.ID, roleID, actorID, now)
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

func (s Service) writeAudit(ctx context.Context, user User, action string, resourceType string, resourceID string, metadata map[string]any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(ctx, audit.Event{
		TenantID:     user.TenantID,
		ActorUserID:  &user.ID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   &resourceID,
		Metadata:     metadata,
	})
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
