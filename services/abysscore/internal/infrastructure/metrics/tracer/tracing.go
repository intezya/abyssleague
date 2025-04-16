package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TraceFn wraps a function with a span and handles error recording.
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

	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	}
}

// TraceVoid runs a function within a span, ignoring errors.
func TraceVoid(ctx context.Context, spanName string, fn func(context.Context)) {
	tracer := otel.Tracer("application")
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	fn(ctx)
}

// TraceValue runs a function within a span and returns a value, ignoring errors.
func TraceValue[T any](ctx context.Context, spanName string, fn func(context.Context) T) T {
	tracer := otel.Tracer("application")
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	return fn(ctx)
}

func TraceValueValue[T any, G any](ctx context.Context, spanName string, fn func(context.Context) (T, G)) (T, G) {
	tracer := otel.Tracer("application")
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	return fn(ctx)
}
