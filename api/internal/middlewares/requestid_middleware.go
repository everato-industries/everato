package middlewares

import (
	"net/http"

	"github.com/dtg-lucifer/everato/api/pkg"
	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := pkg.NewLogger()
		defer logger.Close() // Ensure the logger is closed after the request is processed

		// Generate a unique request ID
		requestId, err := uuid.NewRandom()
		if err != nil {
			logger.StdoutLogger.Error("Error generating UUID for the request")
			logger.FileLogger.Error("Error generating UUID for the request")
			next.ServeHTTP(w, r)
			return
		}

		// Set the request ID in the response header
		w.Header().Set("X-Request-ID", requestId.String())

		// Log the request ID for debugging purposes
		logger.StdoutLogger.Info("Incoming request", "RequestID", requestId.String())
		logger.FileLogger.Info("Incoming request", "RequestID", requestId.String())

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
