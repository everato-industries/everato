package main

import "github.com/dtg-lucifer/everato/api/pkg"

func main() {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the application exits
	logger.StdoutLogger.Info("Seeding")
}
