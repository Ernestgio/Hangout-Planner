package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/cmd"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/db"
)

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
	dbConn, dbCloser, err := db.Connect(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := dbCloser(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize and start the server
	e := cmd.InitializeServer(cfg, dbConn)
	go func() {
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
