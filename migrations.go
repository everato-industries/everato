package main

import (
	"fmt"
	"os"
)

// This function migrates the current database connection to the lates state according to the
// queries and migrations in the migration directory
func MigrateDB() error {
	// Check if the migrations directory exists or not
	// designated location is
	// 	- /internal/db/migrations/*
	_, err := os.Stat("internal/db/migrations")
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exists, checking the embedded fs")
	}

	fs, err := MigrationsFS()
	if err != nil {
		return fmt.Errorf("failed to get migrations filesystem: %w", err)
	}

	return nil
}
