package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
)

type RateLimitStore interface {
	Increment(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error)
}

type RateLimitConfig struct {
	Enabled bool
	Prefix  string
	Limit   int
	Window  time.Duration
	Store   RateLimitStore
	KeyFunc func(c *fiber.Ctx) string
}

type RedisRateLimitStore struct {
	client *redis.Client
}

func NewRedisRateLimitStore(client *redis.Client) RedisRateLimitStore {
	return RedisRateLimitStore{client: client}
}

func (s RedisRateLimitStore) Increment(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}
	if count == 1 {
		if err := s.client.Expire(ctx, key, window).Err(); err != nil {
			return 0, 0, err
		}
		return count, window, nil
	}

	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}
	if ttl < 0 {
		ttl = window
	}
	return count, ttl, nil
}

func RateLimit(cfg RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !cfg.Enabled || cfg.Store == nil || cfg.Limit <= 0 {
			return c.Next()
		}

		window := cfg.Window
		if window <= 0 {
			window = time.Minute
		}

		keyValue := c.IP()
		if cfg.KeyFunc != nil {
			keyValue = cfg.KeyFunc(c)
		}
		key := fmt.Sprintf("rate_limit:%s:%s", cfg.Prefix, hashKey(keyValue))

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		count, retryAfter, err := cfg.Store.Increment(ctx, key, window)
		if err != nil {
			return err
		}

		remaining := max(int64(cfg.Limit)-count, 0)
		c.Set("X-RateLimit-Limit", strconv.Itoa(cfg.Limit))
		c.Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if count > int64(cfg.Limit) {
			seconds := int(retryAfter.Seconds())
			if seconds < 1 {
				seconds = 1
			}
			c.Set(fiber.HeaderRetryAfter, strconv.Itoa(seconds))
			return response.Error(c, fiber.StatusTooManyRequests, "Rate limit exceeded", []map[string]string{
				{"code": "RATE_LIMIT_EXCEEDED"},
			})
		}

		return c.Next()
	}
}

func IPKey(c *fiber.Ctx) string {
	return c.IP()
}

func LoginKey(c *fiber.Ctx) string {
	var body struct {
		Email string `json:"email"`
	}
	_ = json.Unmarshal(c.Body(), &body)
	email := strings.ToLower(strings.TrimSpace(body.Email))
	if email == "" {
		email = "unknown"
	}
	return c.IP() + ":" + email
}

func UserOrIPKey(c *fiber.Ctx) string {
	if claims, ok := authClaims(c); ok {
		return claims
	}
	return c.IP()
}

func authClaims(c *fiber.Ctx) (string, bool) {
	authHeader := c.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		return "", false
	}
	return authHeader, true
}

func hashKey(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
