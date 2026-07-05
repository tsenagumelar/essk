package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Service  ServiceConfig
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Logging  LoggingConfig
	CORS     CORSConfig
}

type ServiceConfig struct {
	Name           string
	Env            string
	Version        string
	DatabaseSchema string
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BodyLimit    int
}

type GRPCConfig struct {
	Port      int
	Upstreams []string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type LoggingConfig struct {
	Level  string
	Pretty bool
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load(serviceName string, defaultHTTPPort int, defaultGRPCPort int, defaultSchema string) Config {
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	postgresPort := getEnv("POSTGRES_PORT", "15432")
	postgresDB := getEnv("POSTGRES_DB", "essk")
	postgresUser := getEnv("POSTGRES_USER", "essk")
	postgresPassword := getEnv("POSTGRES_PASSWORD", "essk")

	return Config{
		Service: ServiceConfig{
			Name:           getEnv("ESSK_SERVICE_NAME", serviceName),
			Env:            getEnv("ESSK_APP_ENV", "local"),
			Version:        getEnv("ESSK_APP_VERSION", "0.1.0"),
			DatabaseSchema: getEnv("ESSK_DB_SCHEMA", defaultSchema),
		},
		HTTP: HTTPConfig{
			Port:         getEnvInt("ESSK_HTTP_PORT", getEnvInt("ESSK_BACKEND_PORT", defaultHTTPPort)),
			ReadTimeout:  getEnvDuration("ESSK_HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvDuration("ESSK_HTTP_WRITE_TIMEOUT", 10*time.Second),
			BodyLimit:    getEnvInt("ESSK_HTTP_BODY_LIMIT", 4*1024*1024),
		},
		GRPC: GRPCConfig{
			Port:      getEnvInt("ESSK_GRPC_PORT", defaultGRPCPort),
			Upstreams: getEnvCSV("ESSK_GRPC_UPSTREAMS", ""),
		},
		Database: DatabaseConfig{
			URL: getEnv(
				"DATABASE_URL",
				fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB),
			),
		},
		Redis: RedisConfig{
			Address:  fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "16379")),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Logging: LoggingConfig{
			Level:  getEnv("ESSK_LOG_LEVEL", "info"),
			Pretty: getEnvBool("ESSK_LOG_PRETTY", true),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ESSK_CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		},
	}
}

func getEnvCSV(key string, fallback string) []string {
	value := getEnv(key, fallback)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
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
