package metrics

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type HandlerMetrics struct {
	counter        metric.Int64Counter
	duration       metric.Float64Histogram
	activeRequests metric.Int64UpDownCounter
}

func NewHandlerMetrics(m *otel.Metrics) *HandlerMetrics {
	return &HandlerMetrics{
		counter:        m.RequestCounter,
		duration:       m.RequestDuration,
		activeRequests: m.ActiveRequests,
	}
}

func (hm *HandlerMetrics) RecordRequest(ctx context.Context, domain string, operation string, method string, status string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("domain", domain),
		attribute.String("operation", operation),
		attribute.String("method", method),
		attribute.String("status", status),
	}

	hm.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	hm.duration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (hm *HandlerMetrics) IncActiveRequests(ctx context.Context, domain string, operation string) {
	attrs := []attribute.KeyValue{
		attribute.String("domain", domain),
		attribute.String("operation", operation),
	}
	hm.activeRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (hm *HandlerMetrics) DecActiveRequests(ctx context.Context, domain string, operation string) {
	attrs := []attribute.KeyValue{
		attribute.String("domain", domain),
		attribute.String("operation", operation),
	}
	hm.activeRequests.Add(ctx, -1, metric.WithAttributes(attrs...))
}
