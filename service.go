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

// InitFunc defines a function type used to initialize a Handle.
// It should perform any setup required before the service runs, such as loading configuration,
// establishing connections, or validating dependencies. If initialization fails, it returns an error.
type InitFunc func(*Handle) error

// RunFunc defines a function type that executes the main logic of the service using a Handle.
// This function should block until the service is stopped or the context is cancelled.
// It returns an error if the service encounters a failure during execution.
type RunFunc func(*Handle) error

// ShutdownFunc defines a function type for gracefully terminating a service.
// It should use the provided Handle to perform cleanup, close resources, and emit diagnostics as needed.
// If shutdown fails, it returns an error.
type ShutdownFunc func(*Handle) error

// New constructs a new Service implementation using the provided name, version, and lifecycle functions.
// The returned Service will use the given InitFunc, RunFunc, and ShutdownFunc for its lifecycle methods.
// If any of the function arguments are nil, the corresponding lifecycle phase will be a no-op.
func New(name, version string, init InitFunc, run RunFunc, shutdown ShutdownFunc) Service {
	return &simpleService{
		name:     name,
		version:  version,
		init:     init,
		run:      run,
		shutdown: shutdown,
	}
}

// simpleService is a basic implementation of the Service interface that delegates
// its lifecycle methods to user-provided function fields.
type simpleService struct {
	name     string
	version  string
	init     func(*Handle) error
	run      func(*Handle) error
	shutdown func(*Handle) error
}

// Name returns the unique identifier name of the service.
func (s simpleService) Name() string {
	return s.name
}

// Version returns the semantic version of the service implementation.
func (s simpleService) Version() string {
	return s.version
}

// Init calls the user-provided initialization function if it is not nil.
// Returns an error if initialization fails, or nil if no function is provided.
func (s simpleService) Init(handle *Handle) error {
	if s.init != nil {
		return s.init(handle)
	}
	return nil
}

// Run calls the user-provided run function if it is not nil.
// Returns an error if the run function fails, or nil if no function is provided.
func (s simpleService) Run(handle *Handle) error {
	if s.run != nil {
		return s.run(handle)
	}
	return nil
}

// Shutdown calls the user-provided shutdown function if it is not nil.
// Returns an error if shutdown fails, or nil if no function is provided.
func (s simpleService) Shutdown(handle *Handle) error {
	if s.shutdown != nil {
		return s.shutdown(handle)
	}
	return nil
}
