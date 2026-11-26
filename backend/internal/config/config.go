// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds every runtime setting of the backend.
type Config struct {
	Env                string
	Port               string
	CORSAllowedOrigins []string
	AutoMigrate        bool
	DB                 DB
}

// DB holds PostgreSQL connection settings.
type DB struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Load reads configuration from the environment. Existing environment
// variables always win over values from a local .env file.
func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		Env:         getString("APP_ENV", "development"),
		Port:        getString("APP_PORT", "8080"),
		AutoMigrate: getBool("DB_AUTO_MIGRATE", true),
		DB: DB{
			Host:     getString("DB_HOST", "localhost"),
			Port:     getString("DB_PORT", "5432"),
			Name:     getString("DB_NAME", ""),
			User:     getString("DB_USER", ""),
			Password: getString("DB_PASSWORD", ""),
			SSLMode:  getString("DB_SSLMODE", "disable"),
		},
	}

	var err error
	if cfg.DB.MaxOpenConns, err = getInt("DB_MAX_OPEN_CONNS", 25); err != nil {
		return Config{}, err
	}
	if cfg.DB.MaxIdleConns, err = getInt("DB_MAX_IDLE_CONNS", 25); err != nil {
		return Config{}, err
	}
	if cfg.DB.ConnMaxLifetime, err = getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute); err != nil {
		return Config{}, err
	}

	for _, origin := range strings.Split(getString("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, origin)
		}
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.DB.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("APP_PORT must be a number: %w", err)
	}
	if _, err := strconv.Atoi(c.DB.Port); err != nil {
		return fmt.Errorf("DB_PORT must be a number: %w", err)
	}
	if len(c.CORSAllowedOrigins) == 0 {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}
	if c.DB.MaxOpenConns < 1 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than zero")
	}
	if c.DB.MaxIdleConns < 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be negative")
	}
	return nil
}

// Addr returns the address Fiber should listen on.
func (c Config) Addr() string {
	return ":" + c.Port
}

// DSN builds a PostgreSQL connection URL from the DB settings.
func (d DB) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   net.JoinHostPort(d.Host, d.Port),
		Path:   "/" + d.Name,
	}
	query := u.Query()
	query.Set("sslmode", d.SSLMode)
	u.RawQuery = query.Encode()
	return u.String()
}

func loadDotEnv() {
	for _, file := range []string{".env", "../.env"} {
		if _, err := os.Stat(file); err == nil {
			// Load never overrides variables that are already set.
			_ = godotenv.Load(file)
			return
		}
	}
}

func getString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func getInt(key string, fallback int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a number: %w", key, err)
	}
	return parsed, nil
}

func getBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration (e.g. 5m): %w", key, err)
	}
	return parsed, nil
}
