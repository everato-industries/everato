//go:build !dev

package main

import (
	"embed"
	"io/fs"
	"net/http"
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
