package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

/*
	For some cases varnamelen linter is disabled because it is infra code.
*/

// TraceFn wraps a function with a span and handles error recording.
//
//nolint:varnamelen
func TraceFn(ctx context.Context, spanName string, fn func(context.Context) error) error {
	tracer := otel.Tracer("application")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	err := fn(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	return err
}

// TraceFnWithResult wraps a function with a span, returns a result and records errors.
//
//nolint:varnamelen
func TraceFnWithResult[T any](ctx context.Context, spanName string, fn func(context.Context) (T, error)) (T, error) {
	tracer := otel.Tracer("application")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	result, err := fn(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	return result, err
}

// AddAttribute adds a custom attribute to the current span if it is recording.
func AddAttribute(ctx context.Context, key string, value interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	switch typedValue := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, typedValue))
	case int:
		span.SetAttributes(attribute.Int(key, typedValue))
	case bool:
		span.SetAttributes(attribute.Bool(key, typedValue))
	}
}

// Trace0 runs a function within a span, ignoring errors.
//
//nolint:varnamelen
func Trace0(ctx context.Context, spanName string, fn func(context.Context)) {
	tracer := otel.Tracer("application")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	fn(ctx)
}

// Trace1 runs a function within a span and returns a result, ignoring errors.
//
//nolint:varnamelen
func Trace1[T any](ctx context.Context, spanName string, fn func(context.Context) T) T {
	tracer := otel.Tracer("application")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	return fn(ctx)
}

// Trace2 runs a function within a span and returns a result, ignoring errors.
//
//nolint:varnamelen
func Trace2[R1, R2 any](ctx context.Context, spanName string, fn func(context.Context) (R1, R2)) (R1, R2) {
	tracer := otel.Tracer("application")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	return fn(ctx)
}
