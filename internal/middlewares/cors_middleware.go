package middlewares

import (
	"net/http"
	"strings"

	"github.com/dtg-lucifer/everato/internal/utils"
)

func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origins_w_comma := utils.GetEnv("CORS_ORIGINS", "localhost:3000,localhost:8080")

		// Always set CORS headers
		host := r.Header.Get("Origin")
		if host == "" || host == "null" {
			// If no Origin header is present, get the host header
			host = r.Header.Get("Host")
		}

		if strings.Contains(origins_w_comma, host) {
			// Set the allowed origins
			w.Header().Set("Access-Control-Allow-Origin", host)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Uncomment the next line if you want to allow credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Uncomment and customize if you want to expose certain headers
		// w.Header().Set("Access-Control-Expose-Headers", "X-Custom-Header")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		h.ServeHTTP(w, r)
	})
}
