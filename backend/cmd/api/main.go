// Command api starts the Fiber HTTP server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/fullstack-starter/backend/internal/config"
	"github.com/example/fullstack-starter/backend/internal/database"
	"github.com/example/fullstack-starter/backend/internal/handler"
	"github.com/example/fullstack-starter/backend/internal/repository"
	"github.com/example/fullstack-starter/backend/internal/service"
	"github.com/gofiber/fiber/v3"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DB, logger)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Warn("failed to close database", "error", err)
		}
	}()

	if cfg.AutoMigrate {
		if err := database.Migrate(cfg.DB, logger); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}
	}

	app := handler.NewApp(cfg, logger)

	api := app.Group("/api")
	api.Get("/health", handler.NewHealthHandler(db, logger).Check)

	itemRepository := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepository)
	handler.NewItemHandler(itemService, logger).Register(api)

	logger.Info("starting API server", "addr", cfg.Addr(), "env", cfg.Env)

	if err := app.Listen(cfg.Addr(), fiber.ListenConfig{
		GracefulContext:       ctx,
		ShutdownTimeout:       shutdownTimeout,
		DisableStartupMessage: cfg.Env == "production",
	}); err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	logger.Info("server stopped gracefully")
	return nil
}
