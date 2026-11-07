// Package db opens the database selected by DB_DRIVER and runs migrations.
package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"

	"boilerplate/api/internal/platform/db/sqlcgen"
)

//go:embed schema.sql
var schema string

// Open connects to sqlite (file path) or postgres (URL) based on driver.
func Open(driver, sqlitePath, databaseURL string) (*sql.DB, error) {
	switch driver {
	case "sqlite":
		var dsn string
		if sqlitePath == ":memory:" {
			dsn = "file::memory:?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
		} else {
			if err := os.MkdirAll(filepath.Dir(sqlitePath), 0o755); err != nil {
				return nil, fmt.Errorf("create sqlite dir: %w", err)
			}
			dsn = "file:" + sqlitePath + "?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
		}
		pool, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("open sqlite: %w", err)
		}
		pool.SetMaxOpenConns(1)
		return pool, nil
	case "postgres":
		pool, err := sql.Open("pgx", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("open postgres: %w", err)
		}
		pool.SetMaxOpenConns(10)
		pool.SetConnMaxLifetime(30 * time.Minute)
		return pool, nil
	default:
		return nil, fmt.Errorf("unknown DB_DRIVER %q", driver)
	}
}

// Migrate executes the idempotent schema. Safe to run on every boot.
func Migrate(ctx context.Context, pool *sql.DB) error {
	if _, err := pool.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

// SeedAdmin creates the admin user when the users table is empty.
func SeedAdmin(ctx context.Context, q *sqlcgen.Queries, username, password string, bcryptCost int) error {
	n, err := q.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("seed count users: %w", err)
	}
	if n > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("seed hash password: %w", err)
	}
	_, err = q.CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           "admin-seed-1",
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("seed create admin: %w", err)
	}
	return nil
}
