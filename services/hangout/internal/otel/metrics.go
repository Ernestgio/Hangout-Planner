package otel

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MeterProvider struct {
	provider *sdkmetric.MeterProvider
}

type Metrics struct {
	// Auth metrics
	AuthCounter  metric.Int64Counter
	AuthDuration metric.Float64Histogram

	// Handler metrics (RED: Rate, Errors, Duration)
	RequestCounter  metric.Int64Counter
	RequestDuration metric.Float64Histogram
	ActiveRequests  metric.Int64UpDownCounter

	// DB operation metrics
	DBOperationDuration metric.Float64Histogram
	DBBatchSize         metric.Int64Histogram

	// gRPC client metrics (calls to file service)
	GRPCClientDuration metric.Float64Histogram
	GRPCClientCounter  metric.Int64Counter
}

func NewMeterProvider(ctx context.Context, cfg *config.OTELConfig) (*MeterProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	var exporter sdkmetric.Exporter

	if cfg.UseStdout {
		exporter, err = stdoutmetric.New()
		if err != nil {
			return nil, err
		}
	} else {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err := grpc.NewClient(
			cfg.Endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		exporter, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, err
		}
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(15*time.Second),
		)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return &MeterProvider{
		provider: mp,
	}, nil
}

func (mp *MeterProvider) Shutdown(ctx context.Context) error {
	if mp.provider != nil {
		return mp.provider.Shutdown(ctx)
	}
	return nil
}

func StartRuntimeMetrics() error {
	return runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second))
}

func InitMetrics() (*Metrics, error) {
	meter := otel.Meter("hangout-service")

	authCounter, err := meter.Int64Counter(
		"hangout.auth.requests",
		metric.WithDescription("Number of authentication requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	authDuration, err := meter.Float64Histogram(
		"hangout.auth.duration",
		metric.WithDescription("Duration of authentication operations"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	requestCounter, err := meter.Int64Counter(
		"hangout.requests.total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"hangout.request.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"hangout.requests.active",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	dbOperationDuration, err := meter.Float64Histogram(
		"hangout.db.operation.duration",
		metric.WithDescription("Database operation duration"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5),
	)
	if err != nil {
		return nil, err
	}

	dbBatchSize, err := meter.Int64Histogram(
		"hangout.db.batch_size",
		metric.WithDescription("Number of items in batch database operations"),
		metric.WithUnit("{item}"),
		metric.WithExplicitBucketBoundaries(1, 5, 10, 25, 50, 100, 250, 500, 1000),
	)
	if err != nil {
		return nil, err
	}

	grpcClientDuration, err := meter.Float64Histogram(
		"hangout.grpc.client.duration",
		metric.WithDescription("gRPC client call duration"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
	)
	if err != nil {
		return nil, err
	}

	grpcClientCounter, err := meter.Int64Counter(
		"hangout.grpc.client.requests",
		metric.WithDescription("Number of gRPC client requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		AuthCounter:         authCounter,
		AuthDuration:        authDuration,
		RequestCounter:      requestCounter,
		RequestDuration:     requestDuration,
		ActiveRequests:      activeRequests,
		DBOperationDuration: dbOperationDuration,
		DBBatchSize:         dbBatchSize,
		GRPCClientDuration:  grpcClientDuration,
		GRPCClientCounter:   grpcClientCounter,
	}, nil
}
