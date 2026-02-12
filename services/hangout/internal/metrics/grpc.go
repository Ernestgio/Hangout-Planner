package metrics

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type GRPCClientMetrics struct {
	counter  metric.Int64Counter
	duration metric.Float64Histogram
}

func NewGRPCClientMetrics(m *otel.Metrics) *GRPCClientMetrics {
	return &GRPCClientMetrics{
		counter:  m.GRPCClientCounter,
		duration: m.GRPCClientDuration,
	}
}

func (gm *GRPCClientMetrics) RecordCall(ctx context.Context, service string, method string, status string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("service", service),
		attribute.String("method", method),
		attribute.String("status", status),
	}

	gm.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	gm.duration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}
