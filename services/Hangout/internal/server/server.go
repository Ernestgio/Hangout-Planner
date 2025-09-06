package server

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/controllers"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/services"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/logging"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func InitializeServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	// 1. Repository Layer
	userRepo := repository.NewUserRepository(db)

	// 2. Service Layer
	userService := services.NewUserService(userRepo)

	// 3. Controller Layer
	responseBuilder := dto.NewStandardResponseBuilder(cfg.Env)
	userController := controllers.NewUserController(userService, responseBuilder)

	// 4. Router Layer
	router := NewRouter(userController)

	// Create a new Echo server instance
	server := echo.New()

	// 5. Initialize Logging
	logging.SetupLogger(server)

	// Register all endpoints using the router
	router.RegisterEndpoints(server)

	return server
}
