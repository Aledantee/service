package service

// Phase represents the lifecycle phase of the service.
type Phase string

const (
	// PhaseWaiting indicates the service is waiting to start.
	PhaseWaiting Phase = "WAITING"
	// PhaseInitializing indicates the service is initializing resources.
	PhaseInitializing Phase = "INITIALIZING"
	// PhaseRunning indicates the service is actively running.
	PhaseRunning Phase = "RUNNING"
	// PhaseStopping indicates the service is in the process of stopping.
	PhaseStopping Phase = "STOPPING"
	// PhaseStopped indicates the service has stopped.
	PhaseStopped Phase = "STOPPED"
	// PhaseError indicates the service has stopped due to an error.
	PhaseError Phase = "ERROR"
)
