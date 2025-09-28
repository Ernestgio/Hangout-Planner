package server

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/validator"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	// Initialize Dependencies / controller layer
	dependencies := InitializeDependencies(cfg, db)

	// Router Layer
	router := NewRouter(dependencies)

	// Create a new Echo server instance
	server := echo.New()

	// validator
	server.Validator = validator.NewValidator()

	// Use middleware
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: constants.LoggerFormat,
	}))

	// Register all endpoints using the router
	router.RegisterEndpoints(server)

	return server
}
