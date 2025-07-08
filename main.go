/*
Everato - Modern Event Management Platform

This is the main entry point for the Everato application.
It initializes configuration, runs database migrations, and starts the HTTP server.
*/
package main

import (
	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/joho/godotenv"
)

// main is the application's entry point. It performs the following operations:
// 1. Loads environment variables from .env file
// 2. Loads application configuration from config.yaml
// 3. Runs database migrations to ensure schema is up-to-date
// 4. Initializes and starts the HTTP server
func main() {
	logger := pkg.NewLogger()
	// Load environment variables from .env file
	// Continues execution even if .env file is not found
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading the .env file, falling back to config.yaml and local env")
	}

	// Load application configuration
	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		logger.Error("Error loading configuration", "err", err.Error())
		panic(err)
	}

	// Print the configuration for debugging purposes
	logger.Info("Configuration loaded successfully:")
	config.PrettyPrint(cfg)

	// Run database migrations to ensure schema is up-to-date
	logger.Info("Running migrations...")
	if err := MigrateDB(cfg); err != nil {
		logger.Error("Error running migrations", "err", err.Error())
		panic(err)
	}
	logger.Info("Migrations completed successfully...")

	// Insert the super users data in the database
	if err := SuperUserInit(cfg); err != nil {
		logger.Error("Error initializing super users", "err", err.Error())
		panic(err)
	}
	logger.Info("Super users initialized successfully...")

	// Initialize and configure the HTTP server
	server := NewServer(cfg)

	// Start the HTTP server - this call is blocking
	if err := server.Start(); err != nil {
		logger.Error("Error starting the server", "err", err.Error())
		panic(err)
	}
}
