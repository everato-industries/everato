//go:build dev
// +build dev

/*
live_fs.go - Development filesystem handling for the Everato application

This file provides live filesystem access for serving static files and
migrations during development. It allows for hot-reloading of files
without rebuilding the application.

Only used in development builds (dev tag). For production, see embedded_fs.go.
*/
package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// PublicFS serves public files (CSS, JS, images) directly from the filesystem in development.
// It returns an HTTP handler that serves static content from the "public" directory,
// allowing developers to see changes without restarting the application.
//
// Returns:
//   - An http.Handler that serves files from the live filesystem
func PublicFS() http.Handler {
	return http.FileServer(http.Dir("public"))
}

// MigrationsFS returns the on-disk migration filesystem for development use.
// It locates the migrations directory and returns an fs.FS interface for accessing
// SQL migration files directly from disk, enabling easy development and testing.
//
// Returns:
//   - An fs.FS filesystem interface to access migration scripts
//   - An error if the migrations directory cannot be found or accessed
func MigrationsFS() (fs.FS, error) {
	// Get the directory of the currently running executable
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("could not determine caller file")
	}
	baseDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(baseDir, "internal", "db", "migrations")

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory not found on disk: %s", migrationsPath)
	}

	return os.DirFS(migrationsPath), nil
}
