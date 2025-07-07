//go:build !dev

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

// PublicFS serves the public files from the embedded filesystem.
// It is used to serve static files like images, CSS, and JavaScript.

//go:embed public
var publicFS embed.FS

func PublicFS() http.Handler {
	fs, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}

	return http.FileServer(http.FS(fs))
}

//go:embed internal/db/migrations
var migrationsFS embed.FS

// Returns the FS with the embedded directory of migrations scripts
// for the postgresql DB
func MigrationsFS() (fs.FS, error) {
	// Check if the migrations directory exists in the embedded filesystem
	_, err := fs.Stat(migrationsFS, "internal/db/migrations")
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("Migrations directory does not exist in the embedded filesystem")
	}

	return migrationsFS, nil
}
