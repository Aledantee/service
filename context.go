package service

import (
	"context"
	"log/slog"

	metricNoop "go.opentelemetry.io/otel/metric/noop"
	traceNoop "go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// loggerKey is a context key for storing a structured logger instance.
// It uses an empty struct to ensure type safety and prevent key collisions.
type loggerKey struct{}

// nameKey is a context key for storing a service name.
// It uses an empty struct to ensure type safety and prevent key collisions.
type nameKey struct{}

// versionKey is a context key for storing a service version.
// It uses an empty struct to ensure type safety and prevent key collisions.
type versionKey struct{}

// otelTracerKey is a context key for storing an OpenTelemetry tracer for a service.
type otelTracerKey struct{}

// otelTracerProviderKey is a context key for storing an OpenTelemetry tracer provider for a service.
type otelTracerProviderKey struct{}

// otelMeterKey is a context key for storing an OpenTelemetry meter for a service.
type otelMeterKey struct{}

// otelMeterProviderKey is a context key for storing an OpenTelemetry meter provider for a service.
type otelMeterProviderKey struct{}

// Logger retrieves a structured logger from the context.
// Returns nil the default logger if no service logger is found in the context.
//
// This function is typically used by services to access their configured
// logger instance for structured logging throughout their lifecycle.
func Logger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return logger
}

// Name retrieves the service name from the context.
// Returns an empty string if no name is found in the context.
//
// This function allows services to access their configured name,
// which is useful for logging, metrics, and service identification.
func Name(ctx context.Context) string {
	name, ok := ctx.Value(nameKey{}).(string)
	if !ok {
		return ""
	}

	return name
}

// Version retrieves the service version from the context.
// Returns an empty string if no version is found in the context.
//
// This function allows services to access their configured version,
// which is useful for logging, metrics, and service identification.
func Version(ctx context.Context) string {
	version, ok := ctx.Value(versionKey{}).(string)
	if !ok {
		return ""
	}

	return version
}

// OTelTracer retrieves an OpenTelemetry service tracer from the context.
// Returns a no-op tracer if no tracer is found in the context.
//
// This function allows services to access their configured tracer instance
// for distributed tracing throughout their lifecycle.
func OTelTracer(ctx context.Context) trace.Tracer {
	tracer, ok := ctx.Value(otelTracerKey{}).(trace.Tracer)
	if !ok {
		return traceNoop.NewTracerProvider().Tracer("noop")
	}

	return tracer
}

// OTelTracerProvider retrieves an OpenTelemetry service tracer provider from the context.
// Returns a no-op tracer provider if no provider is found in the context.
//
// This function allows services to access their configured tracer provider instance
// for creating and managing tracers throughout their lifecycle.
func OTelTracerProvider(ctx context.Context) trace.TracerProvider {
	tracerProvider, ok := ctx.Value(otelTracerProviderKey{}).(trace.TracerProvider)
	if !ok {
		return traceNoop.NewTracerProvider()
	}

	return tracerProvider
}

// OTelMeter retrieves an OpenTelemetry meter from the context.
// Returns a no-op meter if no meter is found in the context.
//
// This function allows services to access their configured meter instance
// for metrics collection throughout their lifecycle.
func OTelMeter(ctx context.Context) metric.Meter {
	meter, ok := ctx.Value(otelMeterKey{}).(metric.Meter)
	if !ok {
		return metricNoop.NewMeterProvider().Meter("noop")
	}

	return meter
}

// OTelMeterProvider retrieves an OpenTelemetry meter provider from the context.
// Returns a no-op meter provider if no provider is found in the context.
//
// This function allows services to access their configured meter provider instance
// for creating and managing meters throughout their lifecycle.
func OTelMeterProvider(ctx context.Context) metric.MeterProvider {
	meterProvider, ok := ctx.Value(otelMeterProviderKey{}).(metric.MeterProvider)
	if !ok {
		return metricNoop.NewMeterProvider()
	}

	return meterProvider
}
