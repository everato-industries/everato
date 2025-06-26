package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	FileLogger   *slog.Logger
	StdoutLogger *slog.Logger
	file         *os.File // Keep a reference to the file
}

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

// Close the file when the logger is no longer needed
func (l *Logger) Close() {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			slog.Error("Failed to close log file", "error", err)
		}
	}
}

// Helper functions
func (l *Logger) Info(msg string, args ...any) {
	l.FileLogger.Info(msg, args...)
	l.StdoutLogger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.FileLogger.Error(msg, args...)
	l.StdoutLogger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.FileLogger.Debug(msg, args...)
	l.StdoutLogger.Debug(msg, args...)
}
