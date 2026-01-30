package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricsRecorder struct {
	metrics *Metrics
}

func NewMetricsRecorder(metrics *Metrics) *MetricsRecorder {
	if metrics == nil {
		return &MetricsRecorder{metrics: nil}
	}
	return &MetricsRecorder{metrics: metrics}
}

func (mr *MetricsRecorder) StartOperation(ctx context.Context, operation string) func(error) {
	if mr == nil || mr.metrics == nil {
		return func(error) {}
	}

	start := time.Now()
	mr.metrics.ActiveRequests.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", operation),
	))

	return func(err error) {
		duration := time.Since(start).Seconds()
		status := "success"
		if err != nil {
			status = "error"
			mr.metrics.RequestErrors.Add(ctx, 1, metric.WithAttributes(
				attribute.String("operation", operation),
			))
		}

		mr.metrics.ActiveRequests.Add(ctx, -1, metric.WithAttributes(
			attribute.String("operation", operation),
		))
		mr.metrics.RequestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("status", status),
		))
		mr.metrics.RequestDuration.Record(ctx, duration, metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("status", status),
		))
	}
}

func (mr *MetricsRecorder) RecordFileSize(ctx context.Context, size int64) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.FileUploadSize.Record(ctx, size)
}

func (mr *MetricsRecorder) RecordS3Operation(ctx context.Context, operation string, duration time.Duration) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.S3OperationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("operation", operation),
	))
}

func (mr *MetricsRecorder) RecordDBOperation(ctx context.Context, operation string, duration time.Duration, batchSize int) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.DBOperationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("operation", operation),
	))
	if batchSize > 0 {
		mr.metrics.DBBatchSize.Record(ctx, int64(batchSize))
	}
}
