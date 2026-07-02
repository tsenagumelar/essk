package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

func runMigrate(args []string) {
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Load()
	m, err := migrate.New("file://migrations", cfg.Database.URL)
	if err != nil {
		fmt.Printf("failed to initialize migrations: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	switch args[0] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("migration up completed")
	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("migration down failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("migration down completed")
	case "version":
		version, dirty, err := m.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			fmt.Printf("migration version failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("version=%d dirty=%t\n", version, dirty)
	default:
		printUsage()
		os.Exit(1)
	}
}
