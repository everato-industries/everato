package main

import (
	"log"

	"github.com/dtg-lucifer/everato/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Print the configuration for debugging purposes
	log.Println("Configuration loaded successfully:")
	config.PrettyPrint(cfg)

	server := NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
