package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	sourceURL := "file://migrations"
	databaseURL := cfg.Database.URL
	command := args[0]
	if args[0] == "shared" {
		if len(args) < 2 {
			printUsage()
			os.Exit(1)
		}
		sourceURL = "file://shared/migrations"
		databaseURL = withMigrationTable(cfg.Database.URL, "shared_schema_migrations")
		command = args[1]
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		fmt.Printf("failed to initialize migrations: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	switch command {
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

func withMigrationTable(databaseURL string, tableName string) string {
	separator := "?"
	if strings.Contains(databaseURL, "?") {
		separator = "&"
	}
	return databaseURL + separator + "x-migrations-table=" + tableName
}
