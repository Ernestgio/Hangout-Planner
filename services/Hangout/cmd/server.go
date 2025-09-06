package cmd

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/logging"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	// 1. Initialize Dependencies / controller layer
	dependencies := InitializeDependencies(cfg, db)

	// 2. Router Layer
	router := NewRouter(dependencies)

	// Create a new Echo server instance
	server := echo.New()

	// 3. Use middleware
	server.Use(middleware.LoggerWithConfig(logging.LoggerConfig()))

	// 4. Register all endpoints using the router
	router.RegisterEndpoints(server)

	return server
}
