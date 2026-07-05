package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	legacyconfig "github.com/tsenagumelar/essk/services/backend/internal/config"
	"github.com/tsenagumelar/essk/services/backend/internal/middleware"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
	authmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/auth"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/product"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
	tenantmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/tenant"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/useradmin"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	appvalidator "github.com/tsenagumelar/essk/services/backend/internal/validator"
	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/handler"
	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/repositories"
	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/usecase"
)

func Register(api fiber.Router, deps service.Dependencies) {
	legacyCfg := legacyconfig.Load()
	rateLimitStore := middleware.NewRedisRateLimitStore(deps.Redis)
	api.Use(middleware.RateLimit(middleware.RateLimitConfig{
		Enabled: legacyCfg.RateLimit.Enabled,
		Prefix:  "global",
		Limit:   legacyCfg.RateLimit.GlobalRPM,
		Window:  legacyCfg.RateLimit.Window,
		Store:   rateLimitStore,
		KeyFunc: middleware.IPKey,
	}))

	repo := repositories.New(deps.DB, deps.Config.Service.DatabaseSchema)
	uc := usecase.New(repo, deps.Config.GRPC.Upstreams)
	h := handler.New(uc)

	group := api.Group("/gateway")
	group.Get("/health", h.Health)
	group.Get("/upstreams", h.Upstreams)

	if deps.DB == nil {
		return
	}

	tokenService := authn.NewTokenService(legacyCfg.Auth)
	rbacRepo := rbac.NewRepository(deps.DB)
	auditRepo := audit.NewRepository(deps.DB)
	auditService := audit.NewService(auditRepo)

	authRepo := authmodule.NewRepository(deps.DB)
	authService := authmodule.NewService(legacyCfg, authRepo, authmodule.NewPasswordHasher(), tokenService).WithAudit(auditService)
	authHandler := authmodule.NewHandler(authService, appvalidator.New())
	authmodule.RegisterRoutes(api, authHandler, tokenService, authmodule.RateLimiters{
		Login: middleware.RateLimit(middleware.RateLimitConfig{
			Enabled: legacyCfg.RateLimit.Enabled,
			Prefix:  "auth_login",
			Limit:   legacyCfg.RateLimit.AuthLoginRPM,
			Window:  legacyCfg.RateLimit.Window,
			Store:   rateLimitStore,
			KeyFunc: middleware.LoginKey,
		}),
		Refresh: middleware.RateLimit(middleware.RateLimitConfig{
			Enabled: legacyCfg.RateLimit.Enabled,
			Prefix:  "auth_refresh",
			Limit:   legacyCfg.RateLimit.AuthRefreshRPM,
			Window:  legacyCfg.RateLimit.Window,
			Store:   rateLimitStore,
			KeyFunc: middleware.UserOrIPKey,
		}),
	})

	auditHandler := audit.NewHandler(auditService)
	audit.RegisterRoutes(api, auditHandler, tokenService, rbacRepo)

	productRepo := product.NewRepository(deps.DB)
	productService := product.NewService(productRepo).WithAudit(auditService)
	productHandler := product.NewHandler(productService, appvalidator.New())
	product.RegisterRoutes(api, productHandler, tokenService, rbacRepo)

	rbacService := rbac.NewService(rbacRepo).WithAudit(auditService)
	rbacHandler := rbac.NewHandler(rbacService, rbacRepo, appvalidator.New())
	rbac.RegisterRoutes(api, rbacHandler, tokenService, rbacRepo)

	userAdminRepo := useradmin.NewRepository(deps.DB)
	userAdminService := useradmin.NewService(userAdminRepo, authmodule.NewPasswordHasher()).WithAudit(auditService)
	userAdminHandler := useradmin.NewHandler(userAdminService, userAdminRepo, appvalidator.New())
	useradmin.RegisterRoutes(api, userAdminHandler, tokenService, rbacRepo)

	tenantRepo := tenantmodule.NewRepository(deps.DB)
	tenantService := tenantmodule.NewService(tenantRepo).WithAudit(auditService)
	tenantHandler := tenantmodule.NewHandler(tenantService, rbacRepo, appvalidator.New())
	tenantmodule.RegisterRoutes(api, tenantHandler, tokenService, rbacRepo)
}
