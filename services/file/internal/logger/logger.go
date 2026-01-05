package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
)

var globalLogger *slog.Logger

func Init(env string, serviceName string) {
	opts := &slog.HandlerOptions{
		Level:     getLogLevel(env),
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	globalLogger = slog.New(handler).With(
		slog.String("service.name", serviceName),
		slog.String("environment", env),
	)

	slog.SetDefault(globalLogger)
}

func getLogLevel(env string) slog.Level {
	switch env {
	case constants.ProductionEnv:
		return slog.LevelInfo
	case constants.DevEnv:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func Logger() *slog.Logger {
	if globalLogger == nil {
		slog.Warn(constants.LoggerNotInitializedWarning)
		Init(constants.DevEnv, "unknown-service")
	}
	return globalLogger
}

// WithContext returns a logger with trace context if available
// When OpenTelemetry is enabled, this will extract trace_id and span_id from context
func WithContext(ctx context.Context) *slog.Logger {
	logger := Logger()

	// TODO: When OTEL is enabled, extract trace info from context
	// import "go.opentelemetry.io/otel/trace"
	// spanCtx := trace.SpanContextFromContext(ctx)
	// if spanCtx.IsValid() {
	// 	logger = logger.With(
	// 		slog.String("trace_id", spanCtx.TraceID().String()),
	// 		slog.String("span_id", spanCtx.SpanID().String()),
	// 	)
	// }

	return logger
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).Info(msg, args...)
}

// Debug logs a debug message with context
func Debug(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).Debug(msg, args...)
}

// Warn logs a warning message with context
func Warn(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).Warn(msg, args...)
}

// Error logs an error message with context and error details
func Error(ctx context.Context, msg string, err error, args ...any) {
	logArgs := append([]any{slog.Any("error", err)}, args...)
	WithContext(ctx).Error(msg, logArgs...)
}

// LogDuration logs the duration of an operation
func LogDuration(ctx context.Context, operation string, startTime time.Time, additionalFields ...any) {
	duration := time.Since(startTime)
	args := append([]any{
		slog.String("operation", operation),
		slog.Duration("duration_ms", duration),
	}, additionalFields...)

	WithContext(ctx).Info("operation completed", args...)
}

// LogGRPCRequest logs gRPC request details
func LogGRPCRequest(ctx context.Context, method string, args ...any) {
	logArgs := append([]any{
		slog.String("grpc.method", method),
		slog.String("type", "grpc_request"),
	}, args...)

	WithContext(ctx).Info("gRPC request received", logArgs...)
}

// LogGRPCResponse logs gRPC response details
func LogGRPCResponse(ctx context.Context, method string, statusCode string, duration time.Duration, args ...any) {
	logArgs := append([]any{
		slog.String("grpc.method", method),
		slog.String("grpc.status", statusCode),
		slog.Duration("duration_ms", duration),
		slog.String("type", "grpc_response"),
	}, args...)

	WithContext(ctx).Info("gRPC request completed", logArgs...)
}
