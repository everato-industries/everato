//go:build !dev
// +build !dev

/*
migrations_prod.go - Production database migration handler for the Everato application

This file manages database schema migrations in production builds. It uses embedded migration
files rather than accessing the filesystem directly, ensuring that all required SQL scripts
are included in the binary itself for reliable deployment.

Only used in production builds (!dev). For development migrations, see migrations_dev.go.
*/
package main

import (
	"database/sql"
	"fmt"

	"github.com/dtg-lucifer/everato/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // Import PostgreSQL driver
)

// MigrateDB applies pending database migrations to bring the schema up to date.
// It uses migration files embedded in the binary for production environments.
//
// This function:
//  1. Gets the embedded migrations filesystem
//  2. Creates a source driver from the embedded migrations
//  3. Establishes a database connection
//  4. Creates a database driver for migrations
//  5. Runs all pending migrations
//
// Parameters:
//   - cfg: Application configuration containing database connection details
//
// Returns:
//   - An error if any part of the migration process fails
func MigrateDB(cfg *config.Config) error {
	// Get the embedded migrations filesystem
	// Migration files are embedded in the binary in production mode
	m_fs, err := MigrationsFS()
	if err != nil {
		return fmt.Errorf("failed to get migrations filesystem: %w", err)
	} else if m_fs == nil {
		return fmt.Errorf("migrations filesystem is nil, ensure migrations directory exists")
	}

	// Create a source driver using the embedded migration files
	// This allows the migrate library to read SQL files from our binary
	source_driver, err := iofs.New(m_fs, "internal/db/migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver for migrations: %w", err)
	}

	// Construct PostgreSQL connection string from configuration
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.Name,
	)

	// Open database connection using the pgx driver
	// The "pgx" driver is imported via blank import above
	db, err := sql.Open("pgx", url)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Create a database driver for the migrate library
	// This handles the actual execution of SQL against the database
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Initialize the migrate instance with our source and database drivers
	m, err := migrate.NewWithInstance("iofs", source_driver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Apply all pending migrations
	// ErrNoChange is acceptable - it just means we're already at the latest version
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
