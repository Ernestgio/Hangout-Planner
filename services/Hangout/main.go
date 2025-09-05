package main

import (
	"log"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server := InitializeServer(cfg)

	server.Logger.Fatal(server.Start(":" + cfg.AppPort))
}
