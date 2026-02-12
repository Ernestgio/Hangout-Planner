package metrics

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type DBMetrics struct {
	duration  metric.Float64Histogram
	batchSize metric.Int64Histogram
}

func NewDBMetrics(m *otel.Metrics) *DBMetrics {
	return &DBMetrics{
		duration:  m.DBOperationDuration,
		batchSize: m.DBBatchSize,
	}
}

func (dm *DBMetrics) RecordOperation(ctx context.Context, operation string, table string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("table", table),
	}

	dm.duration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (dm *DBMetrics) RecordBatchSize(ctx context.Context, operation string, table string, size int) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("table", table),
	}

	dm.batchSize.Record(ctx, int64(size), metric.WithAttributes(attrs...))
}
