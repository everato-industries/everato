package main

import (
	"fmt"
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

	fmt.Printf("Loaded config: %#v\n", cfg)

	server := NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
