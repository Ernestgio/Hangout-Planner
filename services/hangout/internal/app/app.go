package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/db"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/handlers"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/validator"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/router"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type App struct {
	server *echo.Echo
	db     *gorm.DB
	closer func() error
	cfg    *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	dbConn, dbCloser, err := db.Connect(cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(dbConn); err != nil {
		return nil, err
	}

	// Initialize utils
	responseBuilder := response.NewBuilder(cfg.Env == constants.ProductionEnv)
	jwtUtils := utils.NewJWTUtils(cfg.JwtConfig)
	bcryptUtils := utils.NewBcryptUtils(bcrypt.DefaultCost)

	// Repository Layer
	userRepo := repository.NewUserRepository(dbConn)

	// Service Layer
	userService := services.NewUserService(userRepo, bcryptUtils)
	authService := services.NewAuthService(userService, jwtUtils, bcryptUtils)

	// handler Layer
	authhandler := handlers.NewAuthHandler(authService, responseBuilder)

	// Server Setup
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: constants.LoggerFormat}))
	e.Use(middleware.Gzip())

	router.NewRouter(e, authhandler)

	return &App{
		server: e,
		db:     dbConn,
		closer: dbCloser,
		cfg:    cfg,
	}, nil
}

func (a *App) Start() error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- a.server.Start(":" + a.cfg.AppPort)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Received interrupt signal, shutting down...")
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.closer(); err != nil {
		log.Printf("Error closing database connection: %v", err)
		return err
	}

	return nil
}
