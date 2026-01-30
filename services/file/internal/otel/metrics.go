package otel

import (
	"context"
	"time"

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
	// Request metrics (RED: Rate, Errors, Duration)
	RequestCounter  metric.Int64Counter
	RequestErrors   metric.Int64Counter
	RequestDuration metric.Float64Histogram
	ActiveRequests  metric.Int64UpDownCounter

	// File upload metrics
	FileUploadSize metric.Int64Histogram

	// S3 operation metrics
	S3OperationDuration metric.Float64Histogram

	// DB operation metrics
	DBOperationDuration metric.Float64Histogram
	DBBatchSize         metric.Int64Histogram

	// gRPC connection metrics
	GRPCActiveConnections metric.Int64UpDownCounter
	GRPCConnectionEvents  metric.Int64Counter
}

func NewMeterProvider(ctx context.Context, cfg Config) (*MeterProvider, error) {
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

func InitMetrics(meterName string) (*Metrics, error) {
	meter := otel.Meter(meterName)

	requestCounter, err := meter.Int64Counter(
		"file_service.requests.total",
		metric.WithDescription("Total number of file service requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	requestErrors, err := meter.Int64Counter(
		"file_service.requests.errors",
		metric.WithDescription("Total number of failed file service requests"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"file_service.request.duration",
		metric.WithDescription("Duration of file service requests"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"file_service.requests.active",
		metric.WithDescription("Number of active in-flight requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	fileUploadSize, err := meter.Int64Histogram(
		"file_service.upload.size",
		metric.WithDescription("Size of uploaded files in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	s3OperationDuration, err := meter.Float64Histogram(
		"file_service.s3.operation.duration",
		metric.WithDescription("Duration of S3 operations"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	dbOperationDuration, err := meter.Float64Histogram(
		"file_service.db.operation.duration",
		metric.WithDescription("Duration of database operations"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	dbBatchSize, err := meter.Int64Histogram(
		"file_service.db.batch.size",
		metric.WithDescription("Number of records in batch database operations"),
		metric.WithUnit("{record}"),
	)
	if err != nil {
		return nil, err
	}

	grpcActiveConnections, err := meter.Int64UpDownCounter(
		"file_service.grpc.connections.active",
		metric.WithDescription("Number of active gRPC connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	grpcConnectionEvents, err := meter.Int64Counter(
		"file_service.grpc.connection.events",
		metric.WithDescription("gRPC connection lifecycle events"),
		metric.WithUnit("{event}"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		RequestCounter:        requestCounter,
		RequestErrors:         requestErrors,
		RequestDuration:       requestDuration,
		ActiveRequests:        activeRequests,
		FileUploadSize:        fileUploadSize,
		S3OperationDuration:   s3OperationDuration,
		DBOperationDuration:   dbOperationDuration,
		DBBatchSize:           dbBatchSize,
		GRPCActiveConnections: grpcActiveConnections,
		GRPCConnectionEvents:  grpcConnectionEvents,
	}, nil
}

func StartRuntimeInstrumentation() error {
	return runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	)
}
