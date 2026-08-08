package database

import (
	"database/sql"
	"fmt"

	"github.com/dlukt/sqlc-pgx-starter/internal/migrations"
	"github.com/pressly/goose/v3"
	// Register the pgx driver under the name "pgx" for database/sql, which is
	// what goose uses internally.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// withEmbeddedGoose runs fn against a database/sql connection using the
// migration files embedded in the binary.
func withEmbeddedGoose(databaseURL string, fn func(*sql.DB) error) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.Files)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return fn(db)
}

// Migrate applies all pending migrations.
func Migrate(databaseURL string) error {
	return withEmbeddedGoose(databaseURL, func(db *sql.DB) error {
		if err := goose.Up(db, "sql"); err != nil {
			return fmt.Errorf("migrate up: %w", err)
		}
		return nil
	})
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(databaseURL string) error {
	return withEmbeddedGoose(databaseURL, func(db *sql.DB) error {
		if err := goose.Down(db, "sql"); err != nil {
			return fmt.Errorf("migrate down: %w", err)
		}
		return nil
	})
}

// MigrateRedo rolls back and re-applies the most recently applied migration.
func MigrateRedo(databaseURL string) error {
	return withEmbeddedGoose(databaseURL, func(db *sql.DB) error {
		if err := goose.Redo(db, "sql"); err != nil {
			return fmt.Errorf("migrate redo: %w", err)
		}
		return nil
	})
}

// MigrateStatus prints the current migration status to stdout.
func MigrateStatus(databaseURL string) error {
	return withEmbeddedGoose(databaseURL, func(db *sql.DB) error {
		return goose.Status(db, "sql")
	})
}

// CreateMigration writes a new, empty migration pair into the source tree at
// internal/migrations/sql. It writes to the real filesystem (not the embed FS).
func CreateMigration(name string) error {
	goose.SetBaseFS(nil)
	db, err := sql.Open("pgx", "")
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()
	return goose.Create(db, "internal/migrations/sql", name, "sql")
}
