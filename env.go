package service

import (
	"os"
	"strings"
)

// OtelEnabled checks if OpenTelemetry is enabled by checking the OTEL_ENABLED environment variable.
// Returns true if the environment variable is set to any non-empty value.
// Returns false when the variable is unset or set to an empty string.
//
// This provides an opt-in mechanism for OpenTelemetry functionality, as opposed to the
// OpenTelemetry specification's opt-out approach using OTEL_SDK_DISABLED.
func OtelEnabled() bool {
	v, ok := os.LookupEnv("OTEL_ENABLED")
	if ok && v != "" {
		return true
	}

	return false
}

// IsDebugEnabled checks if debug mode is enabled by checking the DEBUG environment variable.
// Returns true if the environment variable is set to any non-empty value.
// Returns false when the variable is unset or set to an empty string.
//
// This provides an opt-in mechanism for debug functionality, allowing services
// to enable additional logging, diagnostics, or development features.
func IsDebugEnabled() bool {
	v, ok := os.LookupEnv("DEBUG")
	if ok && v != "" {
		return true
	}

	return false
}

// LogLevel retrieves the configured log level from the LOG_LEVEL environment variable.
// Returns the value of LOG_LEVEL if set to a non-empty string.
// Returns "info" as the default level when the variable is unset or empty.
//
// This function allows services to configure their logging verbosity through
// environment variables, supporting standard log levels like debug, info, warn, and error.
func LogLevel() string {
	v, ok := os.LookupEnv("LOG_LEVEL")
	if ok && v != "" {
		return v
	}

	return "info"
}

// LogFormat retrieves the configured log format from the LOG_FORMAT environment variable.
// Returns the value of LOG_FORMAT if set to a non-empty string.
// Returns "json" as the default format when the variable is unset or empty.
//
// This function allows services to configure their logging format through
// environment variables, supporting formats like json and text.
func LogFormat() string {
	v, ok := os.LookupEnv("LOG_FORMAT")
	if ok && v != "" {
		return v
	}

	return "json"
}

// IsJsonLogEnabled checks if JSON log format is enabled by checking the LOG_FORMAT environment variable.
// Returns true if the LOG_FORMAT environment variable is set to "json" or "structured".
// Returns false when LOG_FORMAT is unset, empty, or set to any other value.
func IsJsonLogEnabled() bool {
	switch strings.ToLower(LogFormat()) {
	case "json", "structured":
		return true
	}

	return false
}

// IsPrettyLogEnabled checks if pretty text log format is enabled by checking the LOG_FORMAT environment variable.
// Returns true if the LOG_FORMAT environment variable is set to "text" or "pretty".
// Returns false when LOG_FORMAT is unset, empty, or set to any other value.
func IsPrettyLogEnabled() bool {
	switch strings.ToLower(LogFormat()) {
	case "text", "pretty":
		return true
	}

	return false
}
