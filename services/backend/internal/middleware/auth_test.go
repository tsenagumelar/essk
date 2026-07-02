package middleware

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

type fakePermissionStore struct {
	allowed bool
	userID  uuid.UUID
	code    string
}

func (s *fakePermissionStore) UserHasPermission(_ context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	s.userID = userID
	s.code = permissionCode
	return s.allowed, nil
}

func TestRequirePermissionAllowsAuthorizedUser(t *testing.T) {
	tokenService := testTokenService()
	userID := uuid.New()
	token := mustToken(t, tokenService, userID)
	store := &fakePermissionStore{allowed: true}

	app := fiber.New()
	app.Get("/protected", RequireAuth(tokenService), RequirePermission(store, "tenants:read"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	if store.userID != userID {
		t.Fatalf("expected permission check user id %s, got %s", userID, store.userID)
	}
	if store.code != "tenants:read" {
		t.Fatalf("expected permission code tenants:read, got %s", store.code)
	}
}

func TestRequirePermissionRejectsForbiddenUser(t *testing.T) {
	tokenService := testTokenService()
	token := mustToken(t, tokenService, uuid.New())
	store := &fakePermissionStore{allowed: false}

	app := fiber.New()
	app.Get("/protected", RequireAuth(tokenService), RequirePermission(store, "tenants:read"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func testTokenService() authn.TokenService {
	return authn.NewTokenService(config.AuthConfig{
		Issuer:         "essk-test",
		SigningKey:     "test-secret",
		AccessTokenTTL: time.Minute,
	})
}

func mustToken(t *testing.T, tokenService authn.TokenService, userID uuid.UUID) string {
	t.Helper()
	token, _, err := tokenService.CreateAccessToken(userID, nil, "user@essk.local", time.Now().UTC())
	if err != nil {
		t.Fatalf("create token: %v", err)
	}
	return token
}
