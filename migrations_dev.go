//go:build dev
// +build dev

package main

import (
	"database/sql"
	"fmt"

	"github.com/dtg-lucifer/everato/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// This function migrates the current database connection to the lates state according to the
// queries and migrations in the migration directory
func MigrateDB(cfg *config.Config) error {
	// Check if the migrations directory exists or not
	// designated location is
	// 	- /internal/db/migrations/*
	m_fs, err := MigrationsFS()
	if err != nil {
		return fmt.Errorf("failed to get migrations filesystem: %w", err)
	} else if m_fs == nil {
		return fmt.Errorf("migrations filesystem is nil, ensure migrations directory exists")
	}

	source_driver, err := iofs.New(m_fs, ".")
	if err != nil {
		return fmt.Errorf("failed to create source driver for migrations: %w", err)
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.Name,
	)

	db, err := sql.Open("pgx", url)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source_driver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
