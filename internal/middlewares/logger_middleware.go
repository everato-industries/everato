// Package middlewares provides HTTP middleware components for the Everato application.
// These middlewares handle cross-cutting concerns like authentication, logging, and request handling.
package middlewares

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
)

// LoggerMiddleware is a middleware that logs HTTP request information.
//
// This middleware:
// 1. Records the start time of each request
// 2. Passes the request to the next handler in the chain
// 3. Calculates request duration after handler completes
// 4. Logs request details to both STDOUT and a log file (logs/events.log)
//
// Each log entry includes:
// - HTTP method (GET, POST, etc.)
// - Request path
// - Response status code
// - Request duration
// - Client IP address
//
// Parameters:
//   - next: The next handler in the middleware chain
//
// Returns:
//   - An http.Handler that logs requests and calls the next handler
func LoggerMiddleware(next http.Handler) http.Handler {
	// Initialize the logs directory if it doesn't exist
	// This ensures the log directory is available before any requests are processed
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// Create the directory with full permissions
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			slog.Error("Failed to create logs directory", "error", err)
		}
	}

	// Initialize the log file if it doesn't exist
	// This ensures the log file is available before any requests are processed
	if _, err := os.Stat("logs/events.log"); os.IsNotExist(err) {
		// Create the log file with appropriate permissions
		_, err := os.Create("logs/events.log")
		if err != nil {
			slog.Error("Failed to create log file", "error", err)
		}
	}

	// Create the standard output logger once during middleware initialization
	// This improves performance by avoiding logger creation on each request
	std_logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Return the actual middleware handler function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time to calculate request duration
		start := time.Now()

		// Wrap the response writer to capture the status code
		// Default to 200 OK until explicitly changed by handler
		rw := &utils.ResponseWriter{ResponseWriter: w, StatusCode: 200}

		// Call the next handler in the middleware chain with our wrapped response writer
		next.ServeHTTP(rw, r)

		// Open the log file for appending
		// O_CREATE - Create if it doesn't exist
		// O_WRONLY - Open for writing only
		// O_APPEND - Append to the end of file
		file, err := os.OpenFile("logs/events.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("Failed to open log file", "error", err)
		}
		// Ensure the file is closed after logging to prevent resource leaks
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("Failed to close log file", "error", err)
			}
		}()

		// Create a file logger for this specific request
		file_logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

		// Calculate the total request duration
		duration := time.Since(start)

		// Log to standard output with structured data
		std_logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.StatusCode),
			slog.Duration("duration", duration),
			slog.String("ip", utils.GetIP(r)), // Extract client IP with utility function
		)

		// Log identical information to the log file
		file_logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.StatusCode),
			slog.Duration("duration", duration),
			slog.String("ip", utils.GetIP(r)),
		)
	})
}
