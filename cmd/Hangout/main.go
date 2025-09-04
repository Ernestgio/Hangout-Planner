package main

import (
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

	server.GET("/", func(c echo.Context) error {
		return c.String(200, "Hangout Planner API is running!")
	})

	server.Logger.Fatal(server.Start(":" + port))
}
