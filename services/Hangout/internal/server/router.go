package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(server *echo.Echo) {
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world! Hangout Planner API is running!")
	})

	server.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}
