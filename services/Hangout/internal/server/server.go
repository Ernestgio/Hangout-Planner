package server

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/logging"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func InitializeServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	server := echo.New()
	logging.SetupLogger(server)
	RegisterEndpoints(server)

	return server
}
