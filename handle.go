package service

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"sync"
)

// Handle represents a context-aware service handle with diagnostics, tracing, and metrics capabilities.
type Handle struct {
	ctx            context.Context
	logger         *slog.Logger         // Structured logger instance
	name           string               // Service name
	version        string               // Service version
	tracer         trace.Tracer         // OpenTelemetry tracer
	tracerProvider trace.TracerProvider // OpenTelemetry tracer provider
	meter          metric.Meter         // OpenTelemetry meter
	meterProvider  metric.MeterProvider // OpenTelemetry meter provider

	phase    Phase // Phase represents the current operational phase of the service lifecycle.
	phaseMtx sync.RWMutex
}

// Context returns the context associated with the Handle, enabling propagation of cancellation, deadlines, and metadata.
func (h *Handle) Context() context.Context {
	return h.ctx
}

// Name returns the service name.
func (h *Handle) Name() string {
	return h.name
}

// Version returns the service version.
func (h *Handle) Version() string {
	return h.version
}

// Logger returns the configured logger instance.
func (h *Handle) Logger() *slog.Logger {
	return h.logger
}

// OtelTracer returns the configured OpenTelemetry tracer.
func (h *Handle) OtelTracer() trace.Tracer {
	return h.tracer
}

// OtelTracerProvider returns the configured OpenTelemetry tracer provider.
func (h *Handle) OtelTracerProvider() trace.TracerProvider {
	return h.tracerProvider
}

// OtelMeter returns the configured OpenTelemetry meter.
func (h *Handle) OtelMeter() metric.Meter {
	return h.meter
}

// OtelMeterProvider returns the configured OpenTelemetry meter provider.
func (h *Handle) OtelMeterProvider() metric.MeterProvider {
	return h.meterProvider
}

// Phase returns the current operational phase of the service lifecycle.
func (h *Handle) Phase() Phase {
	h.phaseMtx.RLock()
	defer h.phaseMtx.RUnlock()

	return h.phase
}

func (h *Handle) setPhase(phase Phase) {
	h.phaseMtx.Lock()
	defer h.phaseMtx.Unlock()

	h.phase = phase
}
