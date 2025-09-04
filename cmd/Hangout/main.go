package main

import (
	"Hangout/logging"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	if os.Getenv("ENV") != "PROD" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not loaded: %v", err)
		}
	}
}

func main() {
	server := echo.New()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9000"
	}

	logging.SetupLogger(server)
	RegisterEndpoints(server)
	server.Logger.Fatal(server.Start(":" + port))
}
