package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants/logmsg"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/db"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/handlers"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/logger"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/otel"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/services"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/storage"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/validator"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type App struct {
	server         *grpc.Server
	healthServer   *health.Server
	listener       net.Listener
	db             *gorm.DB
	tracerProvider *otel.TracerProvider
	closer         func() error
	cfg            *config.Config
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {

	// DB Connection
	dbConn, dbCloser, err := db.Connect(cfg.DBConfig)
	if err != nil {
		logger.Error(ctx, logmsg.DBConnectionFailed, err,
			slog.String("host", cfg.DBConfig.DBHost),
			slog.String("port", cfg.DBConfig.DBPort),
		)
		return nil, err
	}

	// Storage Layer
	s3Client, err := storage.NewS3Client(ctx, cfg.S3Config)
	if err != nil {
		logger.Error(ctx, logmsg.S3ConnectionFailed, err,
			slog.String("endpoint", cfg.S3Config.Endpoint),
			slog.String("bucket", cfg.S3Config.BucketName),
		)
		_ = dbCloser()
		return nil, err
	}

	// Repository Layer
	repo := repository.NewMemoryFileRepository(dbConn)

	// Initialize validator
	fileValidator := validator.NewFileValidator()

	// Initialize service
	fileService := services.NewFileService(dbConn, repo, s3Client, fileValidator)

	// Initialize handler
	fileHandler := handlers.NewFileHandler(fileService)

	// Initialize OpenTelemetry (if enabled)
	var tracerProvider *otel.TracerProvider
	if cfg.OTELConfig.Enabled {
		otelCfg := otel.Config{
			ServiceName:    cfg.AppName,
			ServiceVersion: cfg.OTELConfig.ServiceVersion,
			Environment:    cfg.Env,
			Endpoint:       cfg.OTELConfig.Endpoint,
			UseStdout:      cfg.OTELConfig.UseStdout,
		}
		tracerProvider, err = otel.NewTracerProvider(ctx, otelCfg)
		if err != nil {
			logger.Error(ctx, logmsg.OTELInitFailed, err)
			_ = dbCloser()
			return nil, err
		}
		logger.Info(ctx, logmsg.OTELInitialized,
			slog.String("endpoint", cfg.OTELConfig.Endpoint),
			slog.Bool("use_stdout", cfg.OTELConfig.UseStdout),
		)
	}

	// Setup network listener
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(ctx, logmsg.NetworkListenerFailed, err, slog.String("addr", addr))
		_ = dbCloser()
		if tracerProvider != nil {
			_ = tracerProvider.Shutdown(ctx)
		}
		return nil, err
	}

	// Setup gRPC Server with mTLS and OTEL
	var grpcOpts []grpc.ServerOption

	// Add mTLS if enabled
	if cfg.MTLSConfig.Enabled {
		tlsConfig, err := loadServerTLSConfig(cfg.MTLSConfig)
		if err != nil {
			logger.Error(ctx, logmsg.MTLSInitFailed, err)
			_ = dbCloser()
			if tracerProvider != nil {
				_ = tracerProvider.Shutdown(ctx)
			}
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
		logger.Info(ctx, logmsg.MTLSInitialized,
			slog.String("cert_file", cfg.MTLSConfig.CertFile),
		)
	}

	// Add OTEL interceptor
	if cfg.OTELConfig.Enabled {
		grpcOpts = append(grpcOpts,
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	// Register file service
	filepb.RegisterFileServiceServer(grpcServer, fileHandler)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("file.v1.FileService", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	return &App{
		server:         grpcServer,
		healthServer:   healthServer,
		listener:       lis,
		db:             dbConn,
		tracerProvider: tracerProvider,
		closer:         dbCloser,
		cfg:            cfg,
	}, nil
}

func (a *App) Start() error {
	ctx := context.Background()

	logger.Info(ctx, logmsg.GRPCServerListening,
		slog.String("addr", a.listener.Addr().String()),
		slog.String("environment", a.cfg.Env),
		slog.String("database", a.cfg.DBConfig.DBName),
		slog.String("s3_bucket", a.cfg.S3Config.BucketName),
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

	// Shutdown tracer provider to flush pending spans
	if a.tracerProvider != nil {
		if err := a.tracerProvider.Shutdown(ctx); err != nil {
			logger.Error(ctx, logmsg.OTELShutdownFailed, err)
		}
	}

	if err := a.closer(); err != nil {
		logger.Error(ctx, logmsg.DBCloseFailed, err)
		return err
	}

	logger.Info(ctx, logmsg.ShutdownComplete)
	return nil
}

// loadServerTLSConfig loads TLS configuration for gRPC server with mTLS
func loadServerTLSConfig(cfg *config.MTLSConfig) (*tls.Config, error) {
	// Load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, apperrors.ErrMTLSCertLoadFailed
	}

	// Load CA certificate for verifying client certificates
	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, apperrors.ErrMTLSCertLoadFailed
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, apperrors.ErrMTLSCertLoadFailed
	}

	// Configure TLS with mutual authentication
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificate
		ClientCAs:    caCertPool,                     // Verify client cert against this CA
		MinVersion:   tls.VersionTLS12,               // Use TLS 1.2 or higher
	}

	return tlsConfig, nil
}
