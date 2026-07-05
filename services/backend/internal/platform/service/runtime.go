package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/tsenagumelar/essk/services/backend/internal/cache"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	platformlogger "github.com/tsenagumelar/essk/services/backend/internal/platform/logger"
	"github.com/tsenagumelar/essk/services/backend/internal/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Dependencies struct {
	Config config.Config
	Log    zerolog.Logger
	DB     *pgxpool.Pool
	Redis  *redis.Client
}

type RouteRegistrar func(api fiber.Router, deps Dependencies)

func Run(ctx context.Context, cfg config.Config, register RouteRegistrar) error {
	log := platformlogger.New(cfg)

	db, err := connectDatabase(ctx, cfg.Database.URL, log)
	if err != nil {
		log.Warn().Err(err).Msg("database connection is not ready")
	}
	if db != nil {
		defer db.Close()
	}
	redisClient := cache.NewRedis(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)
	defer redisClient.Close()

	deps := Dependencies{Config: cfg, Log: log, DB: db, Redis: redisClient}
	httpServer := newHTTPServer(cfg, deps, register)
	grpcServer, listener, err := newGRPCServer(cfg)
	if err != nil {
		return err
	}
	defer grpcServer.Stop()

	errCh := make(chan error, 2)
	go func() {
		address := fmt.Sprintf(":%d", cfg.HTTP.Port)
		log.Info().Str("address", address).Msg("starting http service")
		errCh <- httpServer.Listen(address)
	}()
	go func() {
		log.Info().Str("address", listener.Addr().String()).Msg("starting grpc service")
		errCh <- grpcServer.Serve(listener)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		grpcServer.GracefulStop()
		return httpServer.ShutdownWithContext(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func connectDatabase(ctx context.Context, databaseURL string, log zerolog.Logger) (*pgxpool.Pool, error) {
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(connectCtx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(connectCtx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func newHTTPServer(cfg config.Config, deps Dependencies, register RouteRegistrar) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.Service.Name,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		BodyLimit:    cfg.HTTP.BodyLimit,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			deps.Log.Error().Err(err).Msg("unhandled service error")
			return response.Error(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
		},
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: "GET,POST,PATCH,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Correlation-ID",
	}))
	app.Use(helmet.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.Service.Name})
	})
	app.Get("/ready", func(c *fiber.Ctx) error {
		status := "ok"
		if deps.DB == nil {
			status = "degraded"
		}
		return c.JSON(fiber.Map{"status": status, "service": cfg.Service.Name})
	})
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"service": cfg.Service.Name, "version": cfg.Service.Version, "env": cfg.Service.Env})
	})

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.Service.Name})
	})
	if register != nil {
		register(api, deps)
	}

	return app
}

func newGRPCServer(cfg config.Config) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(cfg.Service.Name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	return server, listener, nil
}
