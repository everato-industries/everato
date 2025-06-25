package main

import (
	"log"

	"github.com/dtg-lucifer/everato/server/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	server := NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
