package main

import (
	"Hangout/logging"

	"github.com/labstack/echo/v4"
)

func InitializeServer(server *echo.Echo) {
	logging.SetupLogger(server)
	RegisterEndpoints(server)
}

func RegisterEndpoints(server *echo.Echo) {
	server.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello world! Hangout Planner API is running!")
	})
}
