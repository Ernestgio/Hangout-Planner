package main

import "github.com/labstack/echo/v4"

func RegisterEndpoints(server *echo.Echo) {
	server.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello world! Hangout Planner API is running!")
	})
}
