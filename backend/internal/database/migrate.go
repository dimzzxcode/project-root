package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/example/fullstack-starter/backend/internal/config"
	"github.com/example/fullstack-starter/backend/migrations"
	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Migrate applies all pending migrations embedded in the binary. It opens its
// own short-lived connection so the application pool is not affected by the
// migration driver closing its handle.
func Migrate(cfg config.DB, logger *slog.Logger) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open migration connection: %w", err)
	}
	defer func() { _ = db.Close() }()

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	before, _, _ := migrator.Version()

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("database schema is up to date", "version", before)
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	after, _, _ := migrator.Version()
	logger.Info("database migrations applied", "from", before, "to", after)
	return nil
}

// MigrateVersion reports the currently applied migration version.
func MigrateVersion(cfg config.DB) (uint, bool, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return 0, false, fmt.Errorf("open migration connection: %w", err)
	}
	defer func() { _ = db.Close() }()

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("create migration driver: %w", err)
	}

	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return 0, false, fmt.Errorf("load embedded migrations: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		return 0, false, fmt.Errorf("create migrator: %w", err)
	}

	version, dirty, err := migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("read migration version: %w", err)
	}
	return version, dirty, nil
}
