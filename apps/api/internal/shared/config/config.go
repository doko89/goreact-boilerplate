// Package config loads application configuration from the environment,
// with optional .env file support for local development.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config is the full runtime configuration of the API.
type Config struct {
	Env           string // dev | prod | test
	Port          string
	DBDriver      string // sqlite | postgres
	SQLitePath    string
	DatabaseURL   string
	SessionSecret string
	SessionTTL    time.Duration
	AdminUser     string
	AdminPass     string
}

// IsProd reports whether the app runs in production mode.
func (c Config) IsProd() bool { return c.Env == "prod" }

// Load reads .env (if found upward from cwd) then resolves Config.
// Explicit environment variables always win over .env values.
func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		Env:           get("APP_ENV", "dev"),
		Port:          get("APP_PORT", "3000"),
		DBDriver:      get("DB_DRIVER", "sqlite"),
		SQLitePath:    get("SQLITE_PATH", "./data/app.db"),
		DatabaseURL:   get("DATABASE_URL", ""),
		SessionSecret: get("SESSION_SECRET", ""),
		AdminUser:     get("ADMIN_USER", "admin"),
		AdminPass:     get("ADMIN_PASS", ""),
	}

	ttlHours, err := strconv.Atoi(get("SESSION_TTL_HOURS", "12"))
	if err != nil || ttlHours <= 0 {
		return Config{}, errors.New("SESSION_TTL_HOURS must be a positive integer")
	}
	cfg.SessionTTL = time.Duration(ttlHours) * time.Hour

	switch cfg.DBDriver {
	case "sqlite", "postgres":
	default:
		return Config{}, fmt.Errorf("DB_DRIVER must be sqlite or postgres, got %q", cfg.DBDriver)
	}
	if cfg.DBDriver == "postgres" && cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required when DB_DRIVER=postgres")
	}
	if cfg.SessionSecret == "" {
		if cfg.IsProd() {
			return Config{}, errors.New("SESSION_SECRET is required in prod")
		}
		cfg.SessionSecret = "dev-only-insecure-secret"
		fmt.Fprintln(os.Stderr, "WARN: using insecure dev SESSION_SECRET, set it in .env")
	}
	if cfg.AdminPass == "" {
		if cfg.IsProd() {
			return Config{}, errors.New("ADMIN_PASS is required in prod (seed admin)")
		}
		cfg.AdminPass = "admin123"
		fmt.Fprintln(os.Stderr, "WARN: using default dev admin password, set ADMIN_PASS in .env")
	}
	return cfg, nil
}

func get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// loadDotEnv parses the nearest .env walking upward from cwd.
// It never overrides variables already present in the environment.
func loadDotEnv() {
	path := findDotEnv()
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, unquote(strings.TrimSpace(val)))
	}
}

func findDotEnv() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for range 6 {
		candidate := filepath.Join(dir, ".env")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
	return ""
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
