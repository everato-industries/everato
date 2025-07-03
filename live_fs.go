//go:build dev

package main

import "net/http"

func PublicFS() http.Handler {
	// In development mode, we serve the public files directly from the local filesystem.
	return http.FileServer(http.Dir("public"))
}
