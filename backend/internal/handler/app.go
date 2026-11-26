package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/example/fullstack-starter/backend/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

const healthCheckTimeout = 2 * time.Second

// NewApp creates the Fiber application with the base middleware configured.
func NewApp(cfg config.Config, logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "fullstack-starter-api",
		BodyLimit:    1 << 20, // 1 MiB
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: NewErrorHandler(logger),
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: cfg.Env != "production"}))
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{fiber.HeaderContentType, fiber.HeaderAccept},
		MaxAge:       3600,
	}))

	return app
}

// HealthHandler reports whether the API and its database are reachable.
type HealthHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(db *sql.DB, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{db: db, logger: logger}
}

type healthData struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Check handles GET /api/health.
func (h *HealthHandler) Check(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), healthCheckTimeout)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		h.logger.Error("health check failed", "error", err)
		return respondError(c, fiber.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "database is unavailable", nil)
	}

	return respond(c, fiber.StatusOK, healthData{Status: "ok", Database: "up"})
}
