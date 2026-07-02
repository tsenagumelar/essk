package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tsenagumelar/essk/services/backend/internal/authn"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
	"github.com/tsenagumelar/essk/services/backend/internal/database"
	authmodule "github.com/tsenagumelar/essk/services/backend/internal/modules/auth"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/rbac"
)

func runSeed(args []string) {
	if len(args) < 1 || args[0] != "admin" {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database.URL)
	if err != nil {
		fmt.Printf("failed to connect database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := authmodule.NewRepository(db)
	rbacRepo := rbac.NewRepository(db)
	service := authmodule.NewService(cfg, repo, authmodule.NewPasswordHasher(), authn.NewTokenService(cfg.Auth)).WithRBAC(rbacRepo)
	if err := service.SeedAdmin(ctx); err != nil {
		fmt.Printf("failed to seed admin: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("admin seed completed: %s\n", cfg.Seed.AdminEmail)
}
