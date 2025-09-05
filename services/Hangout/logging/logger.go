package logging

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupLogger(server *echo.Echo) {
	server.Use(middleware.LoggerWithConfig(LoggerConfig()))
}

func LoggerConfig() middleware.LoggerConfig {
	return middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n",
	}
}
