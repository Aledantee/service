package service

import "context"

// Service represents a long-running component that follows a specific lifecycle: initialization, running, and shutdown phases.
type Service interface {
	// Name returns the unique identifier name of the service.
	Name() string

	// Version returns the semantic version of the service implementation.
	Version() string

	// Phase returns the current operational phase of the service.
	// The phase indicates whether the service is initializing, running, or shutting down.
	Phase() Phase
	// Init performs the service initialization and setup.
	// This method should handle all one-time setup tasks such as:
	// - Loading configuration
	// - Establishing database connections
	// - Initializing resources
	// - Validating dependencies
	//
	// Init is called before Run and should not block indefinitely.
	// Any blocking operations should be moved to the Run method.
	Init(context.Context) error

	// Run executes the main service logic in a blocking manner.
	// This method will be executed concurrently with other services.
	// The service should remain active until the context is cancelled.
	//
	// Typical implementations include:
	// - Event loops
	// - HTTP servers
	// - Background workers
	// - Message consumers
	Run(context.Context)
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
	Shutdown(context.Context) error

	// Wait blocks until the service completes its execution.
	// This method waits for the service to reach a terminal state
	// (PhaseError or PhaseStopped).
	//
	// Returns the error returned by the service's Run method, or nil
	// if the service completed successfully.
	Wait() error
}
