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
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants/logmsg"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/db"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/grpc"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/handlers"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/validator"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/middlewares"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
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
	server       *echo.Echo
	db           *gorm.DB
	fileClient   grpc.FileService
	closer       func() error
	cfg          *config.Config
	tracerCloser func(context.Context) error
	meterCloser  func(context.Context) error
}

func NewApp(ctx context.Context, cfg *config.Config) (app *App, err error) {

	// OTEL Tracer Provider
	otelCfg := otel.Config{
		ServiceName:     cfg.AppName,
		ServiceVersion:  cfg.OTELConfig.ServiceVersion,
		Environment:     cfg.Env,
		Endpoint:        cfg.OTELConfig.TraceEndpoint,
		UseStdout:       cfg.OTELConfig.UseStdout,
		TraceSampleRate: cfg.OTELConfig.TraceSampleRate,
	}
	tracerProvider, err := otel.NewTracerProvider(ctx, otelCfg)
	if err != nil {
		log.Printf(logmsg.OTELTracerProviderInitFailed, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := tracerProvider.Shutdown(ctx); closeErr != nil {
				log.Printf(logmsg.OTELShutdownFailed, closeErr)
			}
		}
	}()

	// OTEL Metrics Provider
	meterProvider, err := otel.NewMeterProvider(ctx, cfg.OTELConfig)
	if err != nil {
		log.Printf(logmsg.OTELMeterProviderInitFailed, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := meterProvider.Shutdown(ctx); closeErr != nil {
				log.Printf(logmsg.OTELShutdownFailed, closeErr)
			}
		}
	}()

	// Start Go runtime metrics collection
	if err := otel.StartRuntimeMetrics(); err != nil {
		log.Printf(logmsg.OTELRuntimeMetricsFailed, err)
		return nil, err
	}

	// Initialize metrics collectors
	otelMetrics, err := otel.InitMetrics()
	if err != nil {
		log.Printf(logmsg.OTELMetricsInitFailed, err)
		return nil, err
	}

	// Create metrics recorder (handles nil checks internally)
	metricsRecorder := otel.NewMetricsRecorder(otelMetrics)

	log.Println(logmsg.OTELInitialized)

	// DB Connection
	dbConn, dbCloser, err := db.Connect(cfg.DBConfig)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := dbCloser(); closeErr != nil {
				log.Printf(logmsg.DBConnectionCloseFailed, closeErr)
			}
		}
	}()

	// gRPC File Service Client
	fileClient, err := grpc.NewFileServiceClient(cfg.GRPCClientConfig)
	if err != nil {
		log.Printf(logmsg.FileServiceClientInitFailed, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := fileClient.Close(); closeErr != nil {
				log.Printf(logmsg.FileServiceClientCloseFailed, closeErr)
			}
		}
	}()
	log.Printf(logmsg.FileServiceClientInitialized, cfg.GRPCClientConfig.FileServiceURL)

	// Initialize utils
	responseBuilder := response.NewBuilder(cfg.Env == constants.ProductionEnv)
	jwtUtils := utils.NewJWTUtils(cfg.JwtConfig)
	bcryptUtils := utils.NewBcryptUtils(bcrypt.DefaultCost)

	// Repository Layer
	userRepo := repository.NewUserRepository(dbConn, metricsRecorder)
	hangoutRepo := repository.NewHangoutRepository(dbConn, metricsRecorder)
	activityRepo := repository.NewActivityRepository(dbConn, metricsRecorder)
	memoryRepo := repository.NewMemoryRepository(dbConn, metricsRecorder)

	// Service Layer
	userService := services.NewUserService(dbConn, userRepo, bcryptUtils, metricsRecorder)
	authService := services.NewAuthService(userService, jwtUtils, bcryptUtils, metricsRecorder)
	hangoutService := services.NewHangoutService(dbConn, hangoutRepo, activityRepo, metricsRecorder)
	activityService := services.NewActivityService(dbConn, activityRepo, metricsRecorder)
	memoryService := services.NewMemoryService(dbConn, memoryRepo, hangoutRepo, fileClient, metricsRecorder)

	// handler Layer
	authHandler := handlers.NewAuthHandler(authService, responseBuilder)
	hangoutHandler := handlers.NewHangoutHandler(hangoutService, responseBuilder)
	activityHandler := handlers.NewActivityHandler(activityService, responseBuilder)
	memoryHandler := handlers.NewMemoryHandler(memoryService, responseBuilder)

	// Server Setup
	e := echo.New()
	e.Validator = validator.NewValidator()

	// middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: constants.LoggerFormat}))
	e.Use(middleware.Decompress())
	e.Use(middlewares.TracingMiddleware(cfg.AppName))
	e.Use(middlewares.MetricsMiddleware(metricsRecorder))

	router.NewRouter(e, cfg, responseBuilder, authHandler, hangoutHandler, activityHandler, memoryHandler)

	return &App{
		server:       e,
		db:           dbConn,
		fileClient:   fileClient,
		closer:       dbCloser,
		cfg:          cfg,
		tracerCloser: tracerProvider.Shutdown,
		meterCloser:  meterProvider.Shutdown,
	}, nil
}

func (a *App) Start() error {
	errChan := make(chan error, 1)
	go func() {
		addr := ":" + a.cfg.AppPort
		errChan <- a.server.Start(addr)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println(logmsg.AppShuttingDown)
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	shutdownTimeout := time.Duration(constants.GracefulShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if a.tracerCloser != nil {
		if err := a.tracerCloser(ctx); err != nil {
			log.Printf(logmsg.OTELShutdownFailed, err)
		}
	}

	if a.meterCloser != nil {
		if err := a.meterCloser(ctx); err != nil {
			log.Printf(logmsg.OTELShutdownFailed, err)
		}
	}

	if a.fileClient != nil {
		if err := a.fileClient.Close(); err != nil {
			log.Printf(logmsg.FileServiceClientCloseFailed, err)
		}
	}

	if err := a.closer(); err != nil {
		log.Printf(logmsg.DBConnectionCloseFailed, err)
		return err
	}

	return nil
}
