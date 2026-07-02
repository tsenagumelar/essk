package app

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/cache"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
	"github.com/tsenagumelar/essk/services/backend/internal/database"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
	authmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/auth"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
	tenantmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/tenant"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
	appvalidator "github.com/tsenagumelar/essk/services/backend/internal/validator"
)

type App struct {
	cfg    config.Config
	log    zerolog.Logger
	server *fiber.App
	db     *pgxpool.Pool
	redis  *redis.Client
}

func New(cfg config.Config, log zerolog.Logger) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database.URL)
	if err != nil {
		log.Warn().Err(err).Msg("database connection is not ready")
	}

	redisClient := cache.NewRedis(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	application := &App{
		cfg:   cfg,
		log:   log,
		db:    db,
		redis: redisClient,
	}

	server := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		BodyLimit:    cfg.HTTP.BodyLimit,
		ErrorHandler: application.errorHandler,
	})

	application.server = server
	application.registerMiddleware()
	application.registerRoutes()

	return application, nil
}

func (a *App) Listen() error {
	address := fmt.Sprintf(":%d", a.cfg.HTTP.Port)
	a.log.Info().Str("address", address).Msg("starting backend")
	if err := a.server.Listen(address); err != nil {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			a.log.Warn().Err(err).Msg("failed to close redis client")
		}
	}
	if a.db != nil {
		a.db.Close()
	}
	return a.server.ShutdownWithContext(ctx)
}

func (a *App) registerMiddleware() {
	a.server.Use(recover.New())
	a.server.Use(requestid.New())
	a.server.Use(a.requestLogger())
	a.server.Use(cors.New(cors.Config{
		AllowOrigins: a.cfg.CORS.AllowedOrigins,
		AllowMethods: "GET,POST,PATCH,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Correlation-ID",
	}))
	a.server.Use(helmet.New())
}

func (a *App) registerRoutes() {
	a.server.Get("/health", a.health)
	a.server.Get("/ready", a.ready)
	a.server.Get("/version", a.version)

	api := a.server.Group("/api/v1")
	api.Get("/health", a.health)
	api.Get("/ready", a.ready)
	api.Get("/version", a.version)

	rateLimitStore := middleware.NewRedisRateLimitStore(a.redis)
	api.Use(middleware.RateLimit(middleware.RateLimitConfig{
		Enabled: a.cfg.RateLimit.Enabled,
		Prefix:  "global",
		Limit:   a.cfg.RateLimit.GlobalRPM,
		Window:  a.cfg.RateLimit.Window,
		Store:   rateLimitStore,
		KeyFunc: middleware.IPKey,
	}))

	if a.db != nil {
		tokenService := authn.NewTokenService(a.cfg.Auth)
		rbacRepo := rbac.NewRepository(a.db)
		auditRepo := audit.NewRepository(a.db)
		auditService := audit.NewService(auditRepo)

		authRepo := authmodule.NewRepository(a.db)
		authService := authmodule.NewService(a.cfg, authRepo, authmodule.NewPasswordHasher(), tokenService).WithAudit(auditService)
		authHandler := authmodule.NewHandler(authService, appvalidator.New())
		authmodule.RegisterRoutes(api, authHandler, tokenService, authmodule.RateLimiters{
			Login: middleware.RateLimit(middleware.RateLimitConfig{
				Enabled: a.cfg.RateLimit.Enabled,
				Prefix:  "auth_login",
				Limit:   a.cfg.RateLimit.AuthLoginRPM,
				Window:  a.cfg.RateLimit.Window,
				Store:   rateLimitStore,
				KeyFunc: middleware.LoginKey,
			}),
			Refresh: middleware.RateLimit(middleware.RateLimitConfig{
				Enabled: a.cfg.RateLimit.Enabled,
				Prefix:  "auth_refresh",
				Limit:   a.cfg.RateLimit.AuthRefreshRPM,
				Window:  a.cfg.RateLimit.Window,
				Store:   rateLimitStore,
				KeyFunc: middleware.UserOrIPKey,
			}),
		})

		auditHandler := audit.NewHandler(auditService)
		audit.RegisterRoutes(api, auditHandler, tokenService, rbacRepo)

		rbacService := rbac.NewService(rbacRepo).WithAudit(auditService)
		rbacHandler := rbac.NewHandler(rbacService, appvalidator.New())
		rbac.RegisterRoutes(api, rbacHandler, tokenService, rbacRepo)

		tenantRepo := tenantmodule.NewRepository(a.db)
		tenantService := tenantmodule.NewService(tenantRepo).WithAudit(auditService)
		tenantHandler := tenantmodule.NewHandler(tenantService, appvalidator.New())
		tenantmodule.RegisterRoutes(api, tenantHandler, tokenService, rbacRepo)
	}
}

func (a *App) errorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		return response.Error(c, appErr.Status, appErr.Message, []map[string]string{
			{"code": appErr.Code},
		})
	}

	var fiberErr *fiber.Error
	if stderrors.As(err, &fiberErr) {
		return response.Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	a.log.Error().Err(err).Str("request_id", c.GetRespHeader(fiber.HeaderXRequestID)).Msg("unhandled request error")
	return response.Error(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
}
