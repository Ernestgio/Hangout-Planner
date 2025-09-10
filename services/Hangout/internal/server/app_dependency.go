package cmd

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/controllers"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AppDependencies struct {
	userController *controllers.UserController
}

func InitializeDependencies(cfg *config.Config, db *gorm.DB) *AppDependencies {
	// 1. Repository Layer
	userRepo := repository.NewUserRepository(db)

	// 2. Service Layer
	userService := services.NewUserService(userRepo, bcrypt.DefaultCost)

	// 3. Controller Layer
	responseBuilder := dto.NewStandardResponseBuilder(cfg.Env)
	userController := controllers.NewUserController(userService, responseBuilder)

	return &AppDependencies{
		userController: userController,
	}
}
