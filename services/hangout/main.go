package main

import (
	"log"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/app"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
)

// @title 			Hangout Planner API
// @version 		1.0
// @description 	API documentation for the Hangout Planner service
// @host 			localhost
// @BasePath 		/rp-api/hangout-service

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT.
func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Create a new application instance
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}

	// 3. Start the application
	if err := app.Start(); err != nil {
		log.Fatalf("Fatal error running application: %v", err)
	}
}
