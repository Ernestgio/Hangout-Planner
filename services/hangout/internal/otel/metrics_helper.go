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

func (mr *MetricsRecorder) RecordAuth(ctx context.Context, operation string, status string, duration time.Duration) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.AuthCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", operation),
		attribute.String("status", status),
	))
	mr.metrics.AuthDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("operation", operation),
		attribute.String("status", status),
	))
}

func (mr *MetricsRecorder) StartRequest(ctx context.Context, domain string, operation string) func(string) {
	if mr == nil || mr.metrics == nil {
		return func(string) {}
	}

	start := time.Now()
	mr.metrics.ActiveRequests.Add(ctx, 1, metric.WithAttributes(
		attribute.String("domain", domain),
		attribute.String("operation", operation),
	))

	return func(status string) {
		duration := time.Since(start).Seconds()
		method := "POST"

		mr.metrics.ActiveRequests.Add(ctx, -1, metric.WithAttributes(
			attribute.String("domain", domain),
			attribute.String("operation", operation),
		))
		mr.metrics.RequestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("domain", domain),
			attribute.String("operation", operation),
			attribute.String("method", method),
			attribute.String("status", status),
		))
		mr.metrics.RequestDuration.Record(ctx, duration, metric.WithAttributes(
			attribute.String("domain", domain),
			attribute.String("operation", operation),
			attribute.String("method", method),
			attribute.String("status", status),
		))
	}
}

func (mr *MetricsRecorder) RecordDBOperation(ctx context.Context, operation string, table string, duration time.Duration, batchSize int) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.DBOperationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("operation", operation),
		attribute.String("table", table),
	))
	if batchSize > 0 {
		mr.metrics.DBBatchSize.Record(ctx, int64(batchSize), metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("table", table),
		))
	}
}

func (mr *MetricsRecorder) RecordGRPCCall(ctx context.Context, service string, method string, status string, duration time.Duration) {
	if mr == nil || mr.metrics == nil {
		return
	}
	mr.metrics.GRPCClientCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("service", service),
		attribute.String("method", method),
		attribute.String("status", status),
	))
	mr.metrics.GRPCClientDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("service", service),
		attribute.String("method", method),
		attribute.String("status", status),
	))
}
