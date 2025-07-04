package service

import (
	"context"
	"sync"

	"github.com/aledantee/ae"
)

// Manager defines the interface for managing service lifecycle operations.
// It provides methods to start, monitor, and gracefully shut down services.
// The Manager coordinates the execution of multiple services concurrently,
// handling their initialization, running, and shutdown phases.
type Manager interface {
	// Run starts a service and returns its unique identifier.
	// The service will be initialized and begin execution in a separate goroutine.
	// The returned ID can be used for subsequent operations like monitoring
	// the service phase or initiating shutdown.
	//
	// The service follows the standard lifecycle: Init -> Run -> Shutdown.
	// If initialization fails, the service will not start and an error is returned.
	// Returns the service ID for future operations and any initialization error.
	Run(ctx context.Context, service Service) (id string, err error)

	// Phase returns the current operational phase of a service by its ID.
	//
	// Returns PhaseUnspecified if the service ID is not found.
	Phase(id string) Phase

	// Shutdown gracefully terminates a specific service by its ID.
	// This method initiates the service's shutdown process by calling
	// its Shutdown method. The context can be used to set a timeout
	// for the shutdown operation.
	//
	// The shutdown process is asynchronous - this method returns immediately
	// after initiating shutdown. Use Wait to block until shutdown completes.
	// Returns an error if the service cannot be shut down.
	// Does nothing if the service is not found.
	Shutdown(ctx context.Context, id string) error

	// ShutdownAll gracefully terminates all managed services.
	// This method initiates the shutdown process for all services
	// by calling their Shutdown method. The context can be used to set
	// a timeout for the shutdown operation.
	ShutdownAll(ctx context.Context) error

	// Wait blocks until a specific service completes its execution.
	// This method waits for the service to reach a terminal state
	// (PhaseError or PhaseStopped).
	//
	// Returns the error returned by the service's Run method, or nil
	// if the service completed successfully.
	Wait(id string) error

	// WaitAll blocks until all managed services complete their execution.
	// This method waits for all services to reach terminal states
	// (PhaseError or PhaseStopped).
	//
	// Returns a joined error from all services or nil if all services completed successfully.
	WaitAll() error
}

// NewManager creates and returns a new Manager instance.
func NewManager() Manager {
	return &manager{
		services: make(map[string]Service),
	}
}

type manager struct {
	services    map[string]Service
	servicesMtx sync.RWMutex
}

func (m *manager) Run(ctx context.Context, service Service) (string, error) {
	id := service.Name() + "@" + service.Version()
	m.servicesMtx.RLock()
	_, ok := m.services[id]
	m.servicesMtx.RUnlock()

	if ok {
		return "", ae.New().
			Attr("name", service.Name()).
			Attr("version", service.Version()).
			Msg("service already running")
	}

	if err := service.Init(ctx); err != nil {
		return "", ae.New().
			Attr("name", service.Name()).
			Attr("version", service.Version()).
			Cause(err).
			Msg("service failed to initialize")
	}

	m.servicesMtx.Lock()
	m.services[id] = service
	m.servicesMtx.Unlock()

	go service.Run(ctx)

	return id, nil
}

func (m *manager) Phase(id string) Phase {
	m.servicesMtx.RLock()
	service, ok := m.services[id]
	m.servicesMtx.RUnlock()

	if !ok {
		return PhaseUnspecified
	}

	return service.Phase()
}

func (m *manager) Shutdown(ctx context.Context, id string) error {
	m.servicesMtx.RLock()
	service, ok := m.services[id]
	m.servicesMtx.RUnlock()

	if !ok {
		return nil
	}

	return service.Shutdown(ctx)
}

func (m *manager) ShutdownAll(ctx context.Context) error {
	m.servicesMtx.Lock()
	defer m.servicesMtx.Unlock()

	var errs []error
	for _, service := range m.services {
		if err := service.Shutdown(ctx); err != nil {
			errs = append(errs, ae.New().
				Attr("name", service.Name()).
				Attr("version", service.Version()).
				Cause(err).
				Msg("service failed to shutdown"),
			)
		}
	}

	return ae.WrapMany("at least one service failed to shut down", errs...)
}

func (m *manager) Wait(id string) error {
	m.servicesMtx.RLock()
	service, ok := m.services[id]
	m.servicesMtx.RUnlock()

	if !ok {
		return nil
	}

	return service.Wait()
}

func (m *manager) WaitAll() error {
	m.servicesMtx.Lock()
	defer m.servicesMtx.Unlock()

	var errs []error
	for _, service := range m.services {
		if err := service.Wait(); err != nil {
			errs = append(errs, ae.New().
				Attr("name", service.Name()).
				Attr("version", service.Version()).
				Cause(err).
				Msg("service exited with an error"),
			)
		}
	}

	return ae.WrapMany("at least one service failed", errs...)
}
