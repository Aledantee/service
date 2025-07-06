package service

import (
	"context"
	"errors"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/aledantee/ae"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "github.com/aledantee/service"

// Service represents a long-running component that follows a specific lifecycle:
// initialization, running, and shutdown phases. It provides a structured approach
// to managing service state and lifecycle operations.
type Service struct {
	// Name is the unique identifier name of the service.
	Name string
	// Version is the semantic version of the service implementation.
	Version string

	// Init performs the service initialization and setup, if present.
	// This method should handle all one-time setup tasks such as:
	// - Loading configuration
	// - Establishing database connections
	// - Initializing resources
	// - Validating dependencies
	//
	// Init is called before Run and should not block indefinitely.
	// Any blocking operations should be moved to the Run method.
	Init func(context.Context) error

	// Run executes the main service logic in a blocking manner.
	// This method will be executed concurrently with other services.
	// The service should remain active until the context is cancelled.
	//
	// Typical implementations include:
	// - Event loops
	// - HTTP servers
	// - Background workers
	// - Message consumers
	Run func(*Service) error

	// Shutdown gracefully terminates the service and performs cleanup operations.
	// This method should handle all cleanup tasks such as:
	// - Closing database connections
	// - Stopping background workers
	// - Releasing resources
	// - Cancelling ongoing operations
	//
	// Shutdown is called after Run completes or when the service is requested
	// to stop. The context can be used to set a timeout for the shutdown operation.
	// The method should not block indefinitely and should return promptly
	// once cleanup is complete or the context is cancelled.
	Shutdown func(context.Context) error

	// phase tracks the current lifecycle phase of the service.
	// The default value of 0 corresponds to PhaseWaiting.
	phase Phase
	// phaseMtx protects concurrent access to the phase field.
	phaseMtx sync.Mutex
	// runCtx is the context used during service execution.
	runCtx context.Context
}

// String returns a string representation of the service in the format "Name@Version".
func (s *Service) String() string {
	return s.Name + "@" + s.Version
}

// Phase returns the current lifecycle phase of the service.
func (s *Service) Phase() Phase {
	return s.phase
}

// IsRunning returns true if the service is in any phase other than PhaseWaiting.
func (s *Service) IsRunning() bool {
	return s.phase != PhaseWaiting
}

// Context returns the service's execution context. If no context has been set,
// it returns context.Background().
func (s *Service) Context() context.Context {
	if s.runCtx == nil {
		return context.Background()
	}

	return s.runCtx
}

// Logger returns a structured logger configured for the service's context.
func (s *Service) Logger() *slog.Logger {
	return Logger(s.Context())
}

// OtelMeter returns the OpenTelemetry meter for the service's context.
func (s *Service) OtelMeter() metric.Meter {
	return OTelMeter(s.Context())
}

// OtelMeterProvider returns the OpenTelemetry meter provider for the service's context.
func (s *Service) OtelMeterProvider() metric.MeterProvider {
	return OTelMeterProvider(s.Context())
}

// OtelTracer returns the OpenTelemetry tracer for the service's context.
func (s *Service) OtelTracer() trace.Tracer {
	return OTelTracer(s.Context())
}

// OtelTracerProvider returns the OpenTelemetry tracer provider for the service's context.
func (s *Service) OtelTracerProvider() trace.TracerProvider {
	return OTelTracerProvider(s.Context())
}

func (s *Service) ExecuteExit() {
	s.ExecuteExitContext(context.Background())
}

func (s *Service) ExecuteExitContext(ctx context.Context) {
	if err := s.execute(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			ae.Print(err)
		}

		os.Exit(ae.ExitCode(err))
	}
}

func (s *Service) Execute() error {
	return s.execute(context.Background())
}

func (s *Service) ExecuteContext(ctx context.Context) error {
	return s.execute(ctx)
}

func (s *Service) execute(ctx context.Context) error {
	errBuilder := ae.New().
		Attr("service.name", s.Name).
		Attr("service.version", s.Version)

	if s.Run == nil {
		return errBuilder.Msg("service has no Run method")
	}

	s.phaseMtx.Lock()
	if s.phase != PhaseWaiting {
		s.phaseMtx.Unlock()
		return errBuilder.Msg("service is already running")
	}
	s.phaseMtx.Unlock()

	s.runCtx = withName(withVersion(ctx, s.Version), s.Name)

	logger := s.initLogger()

	s.setPhase(PhaseInitializing)

	if OtelEnabled() {
		logger.Info("initializing OpenTelemetry")

		var err error
		if err = s.initOtel(ctx); err != nil {
			return errBuilder.Cause(err).
				Msg("failed to initialize OpenTelemetry")
		}
	}

	if s.Init != nil {
		logger.Info("initializing service")

		if err := s.Init(s.runCtx); err != nil {
			return errBuilder.Cause(err).
				Msg("failed to initialize service")
		}
	}

	logger.Info("starting service")
	s.setPhase(PhaseRunning)

	err := s.Run(s)
	var shutdownErr error
	if s.Shutdown != nil {
		logger.Info("shutting down service")
		s.setPhase(PhaseShuttingDown)
		shutdownErr = s.Shutdown(s.runCtx)
	}

	if shutdownErr != nil {
		s.setPhase(PhaseError)
	} else {
		s.setPhase(PhaseStopped)
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		return errBuilder.Cause(err).
			Related(shutdownErr). // will be ignored if nil
			Msg("service failed")
	} else if shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
		return errBuilder.Cause(shutdownErr).
			Msg("service shutdown failed")
	} else {
		return nil
	}
}

// setPhase updates the service's lifecycle phase in a thread-safe manner.
func (s *Service) setPhase(phase Phase) {
	s.phaseMtx.Lock()
	defer s.phaseMtx.Unlock()

	s.phase = phase
}

// initOtel initializes OpenTelemetry for the service using the default environment variables.
// See the OpenTelemetry documentation for more information on how to configure OpenTelemetry using environment variables.
// https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#general-sdk-configuration
func (s *Service) initOtel(ctx context.Context) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(s.Name),
			semconv.ServiceVersion(s.Version),
		),
	)
	if err != nil {
		return ae.New().
			Attr("service.name", s.Name).
			Attr("service.version", s.Version).
			Cause(err).
			Msg("failed to create resource")
	}

	meterProvider := sdkMetric.NewMeterProvider(sdkMetric.WithResource(res))
	s.runCtx = withOTelMeterProvider(s.runCtx, meterProvider)

	tracerProvider := sdkTrace.NewTracerProvider(sdkTrace.WithResource(res))
	s.runCtx = withOTelTracerProvider(s.runCtx, tracerProvider)

	meter := meterProvider.Meter(otelName)
	s.runCtx = withOTelMeter(s.runCtx, meter)

	tracer := tracerProvider.Tracer(otelName)
	s.runCtx = withOTelTracer(s.runCtx, tracer)

	return nil
}

func (s *Service) initLogger() *slog.Logger {
	logLevel := slog.LevelInfo
	switch strings.ToLower(LogLevel()) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	var handler slog.Handler
	if IsDebugEnabled() {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else if IsPrettyLogEnabled() {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler).
		With("service.name", s.Name).
		With("service.version", s.Version)

	s.runCtx = withLogger(s.runCtx, logger)

	return logger
}
