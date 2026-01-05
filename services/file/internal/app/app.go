package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants/logmsg"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/db"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type App struct {
	server       *grpc.Server
	healthServer *health.Server
	listener     net.Listener
	db           *gorm.DB
	closer       func() error
	cfg          *config.Config
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	dbConn, dbCloser, err := db.Connect(cfg.DBConfig)

	if err != nil {
		logger.Error(ctx, logmsg.DBConnectionFailed, err,
			slog.String("host", cfg.DBConfig.DBHost),
			slog.String("port", cfg.DBConfig.DBPort),
		)
		return nil, err
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(ctx, logmsg.NetworkListenerFailed, err, slog.String("addr", addr))
		return nil, err
	}

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("file.v1.FileService", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	return &App{
		server:       grpcServer,
		healthServer: healthServer,
		listener:     lis,
		db:           dbConn,
		closer:       dbCloser,
		cfg:          cfg,
	}, nil
}

func (a *App) Start() error {
	ctx := context.Background()

	logger.Info(ctx, logmsg.GRPCServerListening,
		slog.String("addr", a.listener.Addr().String()),
		slog.String("environment", a.cfg.Env),
		slog.String("database", a.cfg.DBConfig.DBName),
	)

	errChan := make(chan error, 1)
	go func() {
		if err := a.server.Serve(a.listener); err != nil {
			logger.Error(ctx, logmsg.GRPCServerError, err)
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info(ctx, logmsg.ShutdownSignalReceived, slog.String("signal", sig.String()))
	case err := <-errChan:
		if err != nil {
			logger.Error(ctx, logmsg.ServerTerminatedWithError, err)
			return err
		}
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	shutdownTimeout := time.Duration(constants.GracefulShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	logger.Info(ctx, logmsg.ShutdownInitiating)

	a.healthServer.Shutdown()

	shutdownComplete := make(chan struct{})
	go func() {
		a.server.GracefulStop()
		close(shutdownComplete)
	}()

	select {
	case <-shutdownComplete:
	case <-ctx.Done():
		logger.Warn(ctx, logmsg.ShutdownTimeoutExceeded)
		a.server.Stop()
	}

	if err := a.closer(); err != nil {
		logger.Error(ctx, logmsg.DBCloseFailed, err)
		return err
	}

	logger.Info(ctx, logmsg.ShutdownComplete)
	return nil
}
