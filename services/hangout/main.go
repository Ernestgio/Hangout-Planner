package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/db"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/server"
)

// @title 			Hangout Planner API
// @version 		1.0
// @description 	API documentation for the Hangout Planner service
// @host 			localhost:9000
// @BasePath 		/
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := Run(cfg); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func Run(cfg *config.Config) error {
	// Connect to the database
	dbConn, dbCloser, err := db.Connect(cfg.DBConfig)
	if err != nil {
		return err
	}
	defer func() {
		if err := dbCloser(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run database migrations -> TODO: Move to be outside the service
	if err := db.Migrate(dbConn); err != nil {
		return err
	}

	// Initialize and start the server
	e := server.InitializeServer(cfg, dbConn)
	errChan := make(chan error, 1)
	go func() {
		errChan <- e.Start(":" + cfg.AppPort)
	}()

	// Wait for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server with a timeout of 10 seconds
	select {
	case <-quit:
		log.Println("Received interrupt signal, shutting down...")
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}
