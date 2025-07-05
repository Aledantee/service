package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/aledantee/ae"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

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

// Execute runs the service through its complete lifecycle: initialization, running, and shutdown.
// It manages the service state transitions and handles errors appropriately.
//
// The method performs the following steps:
// 1. Validates that the service has a Run method and is not already running
// 2. Initializes OpenTelemetry if enabled
// 3. Calls the Init method if provided
// 4. Executes the Run method
// 5. Calls the Shutdown method if provided
// 6. Sets the final phase based on success or failure
//
// Returns an error if any step fails, with context.Canceled errors being treated as normal termination.
func (s *Service) Execute(ctx context.Context) error {
	errBuilder := ae.New().
		Attr("name", s.Name).
		Attr("version", s.Version)

	if s.Run == nil {
		return errBuilder.Msg("service has no Run method")
	}

	s.phaseMtx.Lock()
	if s.phase != PhaseWaiting {
		s.phaseMtx.Unlock()
		return errBuilder.Msg("service is already running")
	}
	s.phaseMtx.Unlock()

	s.setPhase(PhaseInitializing)

	if OtelEnabled() {
		var err error
		if err = s.initOtel(ctx); err != nil {
			return errBuilder.Cause(err).
				Msg("failed to initialize OpenTelemetry")
		}
	}

	if s.Init != nil {
		if err := s.Init(s.runCtx); err != nil {
			return errBuilder.Cause(err).
				Msg("failed to initialize service")
		}
	}

	s.setPhase(PhaseRunning)

	err := s.Run(s)
	var shutdownErr error
	if s.Shutdown != nil {
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
	return nil
}
