package main

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/logging"

	"github.com/labstack/echo/v4"
)

func InitializeServer(cfg *config.Config) *echo.Echo {
	server := echo.New()
	logging.SetupLogger(server)
	RegisterEndpoints(server)

	return server
}

func RegisterEndpoints(server *echo.Echo) {
	server.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello world! Hangout Planner API is running!")
	})

	server.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}
