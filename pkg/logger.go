// Package pkg provides core utilities and services used throughout the Everato application.
// This includes JWT handling, logging, template utilities, and other shared functionality.
package pkg

import (
	"log/slog"
	"os"
)

// Logger wraps multiple logging outputs to provide unified logging to both console and file.
// It manages file handles and provides convenience methods for common logging operations.
type Logger struct {
	FileLogger   *slog.Logger // Logger that writes to a log file
	StdoutLogger *slog.Logger // Logger that writes to standard output
	file         *os.File     // Reference to the file handle for proper closing
}

// NewLogger creates and initializes a new Logger instance.
// It sets up logging to both a file (logs/app.log) and standard output,
// creating necessary directories and files if they don't exist.
//
// The logger uses structured logging with the following configuration:
// - File logger: JSON format with source information at INFO level
// - Console logger: Text format with source information at INFO level
//
// Returns:
//   - A configured Logger instance, or nil if initialization fails
func NewLogger() *Logger {
	logger := &Logger{}

	// Create the logging folder if it does not exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			slog.Error("Failed to create logs directory", "error", err)
			return nil
		}
	}

	// Create the log file if it does not exist
	if _, err := os.Stat("logs/app.log"); os.IsNotExist(err) {
		_, err := os.Create("logs/app.log")
		if err != nil {
			slog.Error("Failed to create log file", "error", err)
			return nil
		}
	}

	// Open the log file for writing
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return nil
	}

	// Attach the file to the logger
	logger.FileLogger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	logger.StdoutLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	// Keep a reference to the file for closing later
	logger.file = file

	return logger
}

// Close properly closes the log file handle when the logger is no longer needed.
// This should be called with defer after creating a logger to ensure proper resource cleanup.
//
// Example usage:
//
//	logger := pkg.NewLogger()
//	defer logger.Close()
func (l *Logger) Close() {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			slog.Error("Failed to close log file", "error", err)
		}
	}
}

// Helper functions for common logging operations
// Each logs to both file and console outputs

// Info logs a message at INFO level with optional key-value pairs.
// This is suitable for normal operational information.
//
// Parameters:
//   - msg: The main log message
//   - args: Optional key-value pairs for structured logging (must be even count)
func (l *Logger) Info(msg string, args ...any) {
	l.FileLogger.Info(msg, args...)
	l.StdoutLogger.Info(msg, args...)
}

// Error logs a message at ERROR level with optional key-value pairs.
// This is suitable for error conditions that require attention.
//
// Parameters:
//   - msg: The main error message
//   - args: Optional key-value pairs for structured logging (must be even count)
func (l *Logger) Error(msg string, args ...any) {
	l.FileLogger.Error(msg, args...)
	l.StdoutLogger.Error(msg, args...)
}

// Debug logs a message at DEBUG level with optional key-value pairs.
// This is suitable for detailed troubleshooting information.
// Note that logs will only appear if the log level is set to DEBUG or lower.
//
// Parameters:
//   - msg: The main debug message
//   - args: Optional key-value pairs for structured logging (must be even count)
func (l *Logger) Debug(msg string, args ...any) {
	l.FileLogger.Debug(msg, args...)
	l.StdoutLogger.Debug(msg, args...)
}
