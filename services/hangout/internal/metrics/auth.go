package metrics

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type AuthMetrics struct {
	counter  metric.Int64Counter
	duration metric.Float64Histogram
}

func NewAuthMetrics(m *otel.Metrics) *AuthMetrics {
	return &AuthMetrics{
		counter:  m.AuthCounter,
		duration: m.AuthDuration,
	}
}

func (am *AuthMetrics) RecordOperation(ctx context.Context, operation string, status string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("status", status),
	}

	am.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	am.duration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}
