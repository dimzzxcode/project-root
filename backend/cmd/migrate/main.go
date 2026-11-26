// Command migrate applies the embedded database migrations.
//
// Usage:
//
//	go run ./cmd/migrate          # apply all pending migrations
//	go run ./cmd/migrate up       # same as above
//	go run ./cmd/migrate version  # print the current migration version
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/example/fullstack-starter/backend/internal/config"
	"github.com/example/fullstack-starter/backend/internal/database"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		err = database.Migrate(cfg.DB, logger)
	case "version":
		var version uint
		var dirty bool
		version, dirty, err = database.MigrateVersion(cfg.DB)
		if err == nil {
			logger.Info("current migration version", "version", version, "dirty", dirty)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\nusage: migrate [up|version]\n", command)
		os.Exit(1)
	}

	if err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
}
