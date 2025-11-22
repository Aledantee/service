package service

import (
	"context"
	"log/slog"
	"os"
)

// loggerKey is an unexported type used as the key for storing the logger within context.Context.
type loggerKey struct{}

// WithLogger returns a new context derived from ctx that carries the provided slog.Logger.
//
// The logger can later be retrieved with Logger(ctx).
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithChildLogger returns a new context that stores a child logger created
// from the logger present in ctx, with additional attributes provided in attrs.
//
// If no logger is found in ctx, a default logger is used as the parent.
func WithChildLogger(ctx context.Context, attrs ...any) context.Context {
	return context.WithValue(ctx, loggerKey{}, Logger(ctx).With(attrs...))
}

// Logger extracts the slog.Logger from ctx.
//
// If no logger is found in ctx, it returns a default JSON-logging logger
// that outputs to os.Stderr.
func Logger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		return slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}

	return logger
}
