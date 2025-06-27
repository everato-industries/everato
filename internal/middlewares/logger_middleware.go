package middlewares

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
)

// Logger middleware to log the http requests
//
// This will log the requests to the STDOUT as well as
// in the log file in `logs/events.log` file
func LoggerMiddleware(next http.Handler) http.Handler {
	// Initialize the directory and file once during middleware creation
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// create the directory
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			slog.Error("Failed to create logs directory", "error", err)
		}
	}

	if _, err := os.Stat("logs/events.log"); os.IsNotExist(err) {
		// create the log file
		_, err := os.Create("logs/events.log")
		if err != nil {
			slog.Error("Failed to create log file", "error", err)
		}
	}

	// Create loggers once during middleware initialization
	std_logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &utils.ResponseWriter{ResponseWriter: w, StatusCode: 200}

		next.ServeHTTP(rw, r)

		file, err := os.OpenFile("logs/events.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("Failed to open log file", "error", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("Failed to close log file", "error", err)
			}
		}()

		file_logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

		duration := time.Since(start)
		std_logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.StatusCode),
			slog.Duration("duration", duration),
			slog.String("ip", utils.GetIP(r)),
		)
		file_logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.StatusCode),
			slog.Duration("duration", duration),
			slog.String("ip", utils.GetIP(r)),
		)
	})
}
