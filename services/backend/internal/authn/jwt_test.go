package authn

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

func TestTokenServiceCreateAndParseAccessToken(t *testing.T) {
	service := NewTokenService(config.AuthConfig{
		Issuer:         "essk-test",
		SigningKey:     "test-secret",
		AccessTokenTTL: time.Minute,
	})

	userID := uuid.New()
	tenantID := uuid.New()
	now := time.Now().UTC()

	token, expiresAt, err := service.CreateAccessToken(userID, &tenantID, "admin@essk.local", now)
	if err != nil {
		t.Fatalf("create access token: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if !expiresAt.After(now) {
		t.Fatal("expected future expiry")
	}

	claims, err := service.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, claims.UserID)
	}
	if claims.TenantID == nil || *claims.TenantID != tenantID {
		t.Fatalf("expected tenant id %s, got %v", tenantID, claims.TenantID)
	}
	if claims.Email != "admin@essk.local" {
		t.Fatalf("expected email claim, got %s", claims.Email)
	}
}
