package tracer

import (
	"context"

	"github.com/intezya/pkglib/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Config struct {
	Endpoint           string
	ServiceName        string
	ServiceVersion     string
	ServiceEnvironment string
}

func Init(config *Config) func() {
	ctx := context.Background()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(config.Endpoint),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logger.Log.Warnf("creating OTLP trace exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.ServiceEnvironment),
		),
	)
	if err != nil {
		logger.Log.Warnf("creating resource: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	logger.Log.Info("Tracer initialized successfully: ", config.Endpoint)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Log.Warnf("Error shutting down tracer provider: %v", err)
		}
	}
}
