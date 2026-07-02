package app

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

func TestSecurityAndObservabilityHeaders(t *testing.T) {
	cfg := config.Config{
		App: config.AppConfig{
			Name:    "ESSK",
			Env:     "test",
			Version: "test",
		},
		HTTP: config.HTTPConfig{
			Port:      0,
			BodyLimit: 1024 * 1024,
		},
		CORS: config.CORSConfig{
			AllowedOrigins: "http://localhost:3000",
		},
		RateLimit: config.RateLimitConfig{
			Enabled: false,
		},
	}

	application, err := New(cfg, zerolog.Nop())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	req := httptest.NewRequest(fiber.MethodGet, "/health", nil)
	req.Header.Set(correlationIDHeader, "test-correlation-id")
	resp, err := application.server.Test(req)
	if err != nil {
		t.Fatalf("app test: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get(correlationIDHeader) != "test-correlation-id" {
		t.Fatalf("expected correlation id to be propagated, got %q", resp.Header.Get(correlationIDHeader))
	}
	if resp.Header.Get(fiber.HeaderXRequestID) == "" {
		t.Fatal("expected request id header")
	}
	if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected nosniff header, got %q", resp.Header.Get("X-Content-Type-Options"))
	}
	if resp.Header.Get("X-Frame-Options") == "" {
		t.Fatal("expected X-Frame-Options header")
	}
}
