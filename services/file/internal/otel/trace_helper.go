package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "file-service"

type SpanRecorder struct {
	span trace.Span
}

func StartServiceSpan(ctx context.Context, operation string, attrs ...attribute.KeyValue) (context.Context, *SpanRecorder) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, "service."+operation)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	return ctx, &SpanRecorder{span: span}
}

func StartRepositorySpan(ctx context.Context, operation string, attrs ...attribute.KeyValue) (context.Context, *SpanRecorder) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, "repository."+operation)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	return ctx, &SpanRecorder{span: span}
}

func (sr *SpanRecorder) End() {
	if sr != nil && sr.span != nil {
		sr.span.End()
	}
}

func (sr *SpanRecorder) SetAttributes(attrs ...attribute.KeyValue) {
	if sr != nil && sr.span != nil {
		sr.span.SetAttributes(attrs...)
	}
}

func (sr *SpanRecorder) RecordError(err error) {
	if sr != nil && sr.span != nil && err != nil {
		sr.span.RecordError(err)
		sr.span.SetStatus(codes.Error, err.Error())
	}
}

func (sr *SpanRecorder) SetStatusOk() {
	if sr != nil && sr.span != nil {
		sr.span.SetStatus(codes.Ok, "")
	}
}

func (sr *SpanRecorder) RecordErrorWithStatus(err error) error {
	if err != nil {
		sr.RecordError(err)
	}
	return err
}

func (sr *SpanRecorder) EndWithError(err error) {
	if err != nil {
		sr.RecordError(err)
	} else {
		sr.SetStatusOk()
	}
	sr.End()
}
