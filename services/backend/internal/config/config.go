package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Version string
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BodyLimit    int
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type AuthConfig struct {
	Issuer          string
	SigningKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

type LoggingConfig struct {
	Level  string
	Pretty bool
}

func Load() Config {
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	postgresDB := getEnv("POSTGRES_DB", "essk")
	postgresUser := getEnv("POSTGRES_USER", "essk")
	postgresPassword := getEnv("POSTGRES_PASSWORD", "essk")

	return Config{
		App: AppConfig{
			Name:    getEnv("ESSK_APP_NAME", "ESSK"),
			Env:     getEnv("ESSK_APP_ENV", "local"),
			Version: getEnv("ESSK_APP_VERSION", "0.1.0"),
		},
		HTTP: HTTPConfig{
			Port:         getEnvInt("ESSK_BACKEND_PORT", 8080),
			ReadTimeout:  getEnvDuration("ESSK_HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvDuration("ESSK_HTTP_WRITE_TIMEOUT", 10*time.Second),
			BodyLimit:    getEnvInt("ESSK_HTTP_BODY_LIMIT", 4*1024*1024),
		},
		Database: DatabaseConfig{
			URL: getEnv(
				"DATABASE_URL",
				fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB),
			),
		},
		Redis: RedisConfig{
			Address:  fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			Issuer:          getEnv("JWT_ISSUER", "essk"),
			SigningKey:      getEnv("JWT_SIGNING_KEY", "change-me-local-only"),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TOKEN_TTL", 168*time.Hour),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ESSK_CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("ESSK_LOG_LEVEL", "info"),
			Pretty: getEnvBool("ESSK_LOG_PRETTY", true),
		},
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
