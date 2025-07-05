package service

// Phase represents the operational state of a service.
// It tracks the lifecycle progression from waiting to stopped.
//
//go:generate go tool stringer -type Phase -trimprefix Phase
type Phase int

const (
	// PhaseWaiting indicates the service is waiting to begin initialization.
	// This is the initial state before any lifecycle methods are called.
	PhaseWaiting Phase = iota
	// PhaseInitializing indicates the service is currently executing its Init method.
	// During this phase, the service performs one-time setup tasks.
	PhaseInitializing
	// PhaseRunning indicates the service has completed initialization and is actively
	// executing its Run method. This is the normal operational state.
	PhaseRunning
	// PhaseShuttingDown indicates the service is in the process of shutting down.
	// During this phase, the service performs cleanup operations.
	PhaseShuttingDown
	// PhaseError indicates the service encountered an error during initialization or
	// execution. This is a terminal state that prevents further lifecycle operations.
	PhaseError
	// PhaseStopped indicates the service has completed its Shutdown method and
	// is no longer operational. This is the final state.
	PhaseStopped
)
