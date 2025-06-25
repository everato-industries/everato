package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	FileLogger   *slog.Logger
	StdoutLogger *slog.Logger
}

func NewLogger() *Logger {
	logger := &Logger{}

	// Create a logger that writes to a file
	//
	// Create the logging folder if that is not present the project root
	// `/logs/`
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// create the directory
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			slog.Error("Failed to create logs directory", "error", err)
			return nil
		}
	}

	if _, err := os.Stat("logs/app.log"); os.IsNotExist(err) {
		// create the log file
		_, err := os.Create("logs/app.log")
		if err != nil {
			slog.Error("Failed to create log file", "error", err)
			return nil
		}
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return nil
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close log file", "error", err)
		}
	}()

	// Attach that to the logger
	logger.FileLogger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	logger.StdoutLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	return logger
}

// Helper functions
// But these are not recommended to use because these will mess up the source
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
