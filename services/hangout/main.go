package main

import (
	"context"
	"log"
	"os"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/app"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants/logmsg"
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
	ctx := context.Background()

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf(logmsg.ConfigLoadFailed, err)
		os.Exit(1)
	}

	// 2. Create a new application instance
	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Printf(logmsg.AppCreateFailed, err)
		os.Exit(1)
	}

	// 3. Start the application
	if err := app.Start(); err != nil {
		log.Printf(logmsg.AppTerminatedWithError, err)
		os.Exit(1)
	}

	log.Println(logmsg.AppExitSuccess)
}
