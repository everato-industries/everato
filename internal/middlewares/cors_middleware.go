package middlewares

import "net/http"

func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Uncomment the next line if you want to allow credentials
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

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
