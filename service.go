package service

// Service defines the interface for a long-running component that follows a specific lifecycle
// with initialization, running, and shutdown phases.
type Service interface {
	// Name returns the unique identifier name of the service
	Name() string
	// Version returns the semantic version of the service implementation
	Version() string
	// Init performs one-time service initialization and setup.
	// This should handle tasks like loading config, establishing connections,
	// and validating dependencies. Should not block indefinitely.
	Init(*Handle) error
	// Run executes the main service logic in a blocking manner.
	// The service should remain active until the context is cancelled.
	// This may be executed concurrently with other services.
	Run(*Handle) error
	// Shutdown gracefully terminates the service and cleans up resources.
	// This handles tasks like closing connections and stopping workers.
	// Should not block indefinitely and return once cleanup is complete
	// or context is cancelled.
	Shutdown(*Handle) error
}
