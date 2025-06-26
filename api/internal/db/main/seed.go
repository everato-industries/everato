package main

import (
	"github.com/dtg-lucifer/everato/server/pkg/logger"
)

func main() {
	logger := logger.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the application exits
	logger.StdoutLogger.Info("Seeding")
}
