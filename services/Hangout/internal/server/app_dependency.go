package cmd

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/controllers"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AppDependencies struct {
	authController *controllers.AuthController
}

func InitializeDependencies(cfg *config.Config, db *gorm.DB) *AppDependencies {
	// Initialize utils
	responseBuilder := dto.NewStandardResponseBuilder(cfg.Env)
	jwtUtils := utils.NewJWTUtils(cfg.JWTSecret, cfg.JWTExpirationHours)
	bcryptUtils := utils.NewBcryptUtils(bcrypt.DefaultCost)

	// 1. Repository Layer
	userRepo := repository.NewUserRepository(db)

	// 2. Service Layer
	userService := services.NewUserService(userRepo, bcryptUtils)
	authService := services.NewAuthService(userService, jwtUtils, bcryptUtils)

	// 3. Controller Layer

	authController := controllers.NewAuthController(authService, responseBuilder)

	return &AppDependencies{
		authController: authController,
	}
}
