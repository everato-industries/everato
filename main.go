/*
Everato - Modern Event Management Platform

This is the main entry point for the Everato application.
It initializes configuration, runs database migrations, and starts the HTTP server.
*/
package main

import (
	"log"

	"github.com/dtg-lucifer/everato/config"
	"github.com/joho/godotenv"
)

// main is the application's entry point. It performs the following operations:
// 1. Loads environment variables from .env file
// 2. Loads application configuration from config.yaml
// 3. Runs database migrations to ensure schema is up-to-date
// 4. Initializes and starts the HTTP server
func main() {
	// Load environment variables from .env file
	// Continues execution even if .env file is not found
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Load application configuration
	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Print the configuration for debugging purposes
	log.Println("Configuration loaded successfully:")
	config.PrettyPrint(cfg)

	// Run database migrations to ensure schema is up-to-date
	log.Println("Running migrations...")
	if err := MigrateDB(cfg); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	log.Println("Migrations completed successfully...")

	// Initialize and configure the HTTP server
	server := NewServer(cfg)

	// Start the HTTP server - this call is blocking
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
