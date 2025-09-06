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
	if err := Run(); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func Run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelDB()
	dbConn, dbCloser, err := db.Connect(ctxDB, cfg)
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

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	//Perform a graceful server shutdown with a timeout
	ctxServer, cancelServer := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelServer()

	if err := e.Shutdown(ctxServer); err != nil {
		return err
	}

	return nil
}
