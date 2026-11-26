// Package database opens and configures the PostgreSQL connection pool.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/example/fullstack-starter/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
)

const (
	connectAttempts = 10
	connectInterval = 2 * time.Second
	pingTimeout     = 5 * time.Second
)

// Connect opens a PostgreSQL pool and waits until the database accepts
// connections, which makes startup resilient to a database that is still
// booting.
func Connect(ctx context.Context, cfg config.DB, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	var pingErr error
	for attempt := 1; attempt <= connectAttempts; attempt++ {
		pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
		pingErr = db.PingContext(pingCtx)
		cancel()

		if pingErr == nil {
			logger.Info("connected to database", "host", cfg.Host, "database", cfg.Name)
			return db, nil
		}

		logger.Warn("database is not ready yet", "attempt", attempt, "max_attempts", connectAttempts, "error", pingErr)

		select {
		case <-ctx.Done():
			_ = db.Close()
			return nil, ctx.Err()
		case <-time.After(connectInterval):
		}
	}

	_ = db.Close()
	return nil, fmt.Errorf("database unreachable after %d attempts: %w", connectAttempts, pingErr)
}
