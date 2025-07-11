package service

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/aledantee/ae"
	"github.com/aledantee/ae/errors"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// ExecuteExit runs the service using context.Background() and exits the program
// with the appropriate exit code if an error occurs. If the error is not context.Canceled,
// it prints the error before exiting.
func ExecuteExit(svc Service) {
	ExecuteExitContext(context.Background(), svc)
}

// ExecuteExitContext runs the service with the provided context and exits the program
// with the appropriate exit code if an error occurs. If the error is not context.Canceled,
// it prints the error before exiting. If the error is context.Canceled, it will not print the error.
func ExecuteExitContext(ctx context.Context, svc Service) {
	if err := execute(ctx, svc); err != nil {
		if !errors.Is(err, context.Canceled) {
			ae.Print(err)
		}

		os.Exit(ae.ExitCode(err))
	}
}

// Execute runs the service using context.Background() and returns any error that occurs.
func Execute(svc Service) error {
	return execute(context.Background(), svc)
}

// ExecuteContext runs the service with the provided context and returns any error that occurs.
func ExecuteContext(ctx context.Context, svc Service) error {
	return execute(ctx, svc)
}

// execute contains the internal common execution logic for all Execute* methods.
// It initializes the logger, OpenTelemetry (if enabled), and the service itself,
// then runs and shuts down the service, handling errors appropriately.
func execute(ctx context.Context, svc Service) error {
	errBuilder := ae.New().
		Attr("service.name", svc.Name()).
		Attr("service.version", svc.Version())

	ctx = withName(ctx, svc.Name())
	ctx = withVersion(ctx, svc.Version())

	handle := &Handle{
		ctx:     ctx,
		name:    svc.Name(),
		version: svc.Version(),
	}

	handle.setPhase(PhaseInitializing)

	if err := initLogger(handle); err != nil {
		return errBuilder.Cause(err).
			Msg("failed to initialize logger")
	}

	if OtelEnabled() {
		handle.Logger().Info("initializing OpenTelemetry")

		var err error
		if err = initOtel(handle); err != nil {
			return errBuilder.Cause(err).
				Msg("failed to initialize OpenTelemetry")
		}
	}

	if err := svc.Init(handle); err != nil {
		return errBuilder.Cause(err).
			Msg("failed to initialize service")
	}

	handle.Logger().Info("starting service")
	handle.setPhase(PhaseRunning)

	err := svc.Run(handle)
	if err != nil {
		handle.Logger().Error("service exited with error", "error", err)
	}

	handle.Logger().Info("shutting down service")
	handle.setPhase(PhaseShuttingDown)
	shutdownErr := svc.Shutdown(handle)

	if shutdownErr != nil {
		handle.setPhase(PhaseError)
	} else {
		handle.setPhase(PhaseStopped)
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		return errBuilder.Cause(err).
			Related(shutdownErr). // will be ignored if nil
			Msg("service exited with error")
	} else if shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
		return errBuilder.Cause(shutdownErr).
			Msg("service shutdown failed")
	} else {
		return nil
	}
}

// initOtel configures OpenTelemetry for the service using environment variables.
// It sets up the OpenTelemetry resource, meter provider, tracer provider, meter, and tracer
// for the given service handle. The context is updated with the new providers and returned via the handle.
// For environment variable configuration details, see:
// https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#general-sdk-configuration
func initOtel(handle *Handle) error {
	res, err := resource.New(handle.Context(),
		resource.WithAttributes(
			semconv.ServiceName(handle.Name()),
			semconv.ServiceVersion(handle.Version()),
		),
	)
	if err != nil {
		return ae.New().
			Attr("service.name", handle.Name()).
			Attr("service.version", handle.Version()).
			Cause(err).
			Msg("failed to create resource")
	}

	handle.meterProvider = sdkMetric.NewMeterProvider(sdkMetric.WithResource(res))
	handle.ctx = withOTelMeterProvider(handle.Context(), handle.meterProvider)

	handle.tracerProvider = sdkTrace.NewTracerProvider(sdkTrace.WithResource(res))
	handle.ctx = withOTelTracerProvider(handle.Context(), handle.tracerProvider)

	handle.meter = handle.meterProvider.Meter(otelName)
	handle.ctx = withOTelMeter(handle.Context(), handle.meter)

	handle.tracer = handle.tracerProvider.Tracer(otelName)
	handle.ctx = withOTelTracer(handle.Context(), handle.tracer)

	return nil
}

// initLogger configures the logger for the service handle based on the current log level and formatting options.
// It supports debug, pretty (text), and JSON logging modes. The logger is attached to the handle and
// the context is updated accordingly.
func initLogger(handle *Handle) error {
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

	handle.logger = slog.New(handler).
		With("service.name", handle.Name()).
		With("service.version", handle.Version())
	handle.ctx = withLogger(handle.ctx, handle.logger)

	return nil
}
