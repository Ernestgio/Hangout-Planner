package cmd

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/controllers"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AppDependencies struct {
	authController *controllers.AuthController
}

func InitializeDependencies(cfg *config.Config, db *gorm.DB) *AppDependencies {
	// 1. Repository Layer
	userRepo := repository.NewUserRepository(db)

	// 2. Service Layer
	userService := services.NewUserService(userRepo, bcrypt.DefaultCost)
	authService := services.NewAuthService(userService)

	// 3. Controller Layer
	responseBuilder := dto.NewStandardResponseBuilder(cfg.Env)
	authController := controllers.NewAuthController(authService, responseBuilder)

	return &AppDependencies{
		authController: authController,
	}
}
