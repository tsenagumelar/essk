package middleware

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type memoryRateLimitStore struct {
	mu      sync.Mutex
	counts  map[string]int64
	expires map[string]time.Time
}

func newMemoryRateLimitStore() *memoryRateLimitStore {
	return &memoryRateLimitStore{
		counts:  make(map[string]int64),
		expires: make(map[string]time.Time),
	}
}

func (s *memoryRateLimitStore) Increment(_ context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if expiresAt, ok := s.expires[key]; !ok || now.After(expiresAt) {
		s.counts[key] = 0
		s.expires[key] = now.Add(window)
	}

	s.counts[key]++
	return s.counts[key], time.Until(s.expires[key]), nil
}

func TestRateLimitAllowsRequestsWithinLimit(t *testing.T) {
	app := fiber.New()
	app.Use(RateLimit(RateLimitConfig{
		Enabled: true,
		Prefix:  "test",
		Limit:   2,
		Window:  time.Minute,
		Store:   newMemoryRateLimitStore(),
		KeyFunc: IPKey,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	for i := 0; i < 2; i++ {
		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
		if err != nil {
			t.Fatalf("app test: %v", err)
		}
		if resp.StatusCode != fiber.StatusNoContent {
			t.Fatalf("expected 204, got %d", resp.StatusCode)
		}
	}
}

func TestRateLimitRejectsRequestsOverLimit(t *testing.T) {
	app := fiber.New()
	app.Use(RateLimit(RateLimitConfig{
		Enabled: true,
		Prefix:  "test",
		Limit:   1,
		Window:  time.Minute,
		Store:   newMemoryRateLimitStore(),
		KeyFunc: IPKey,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	first, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("first request: %v", err)
	}
	if first.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected first status 204, got %d", first.StatusCode)
	}

	second, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("second request: %v", err)
	}
	if second.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected second status 429, got %d", second.StatusCode)
	}
	if second.Header.Get(fiber.HeaderRetryAfter) == "" {
		t.Fatal("expected Retry-After header")
	}
}

func TestLoginKeyIncludesEmail(t *testing.T) {
	app := fiber.New()
	var key string
	app.Post("/", func(c *fiber.Ctx) error {
		key = LoginKey(c)
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(`{"email":"Admin@ESSK.Local"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	if _, err := app.Test(req); err != nil {
		t.Fatalf("app test: %v", err)
	}
	if !strings.Contains(key, "admin@essk.local") {
		t.Fatalf("expected key to include normalized email, got %s", key)
	}
}
