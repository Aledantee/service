package service

import "os"

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

// IsJsonLogEnabled checks if JSON logging is enabled by checking the LOG_JSON environment variable.
// Returns true if the environment variable is set to any non-empty value.
// Returns false when the variable is unset or set to an empty string.
//
// This provides an opt-in mechanism for structured JSON logging output,
// which is useful for log aggregation systems and machine-readable log processing.
func IsJsonLogEnabled() bool {
	v, ok := os.LookupEnv("LOG_JSON")
	if ok && v != "" {
		return true
	}

	return false
}

// IsPrettyLogEnabled checks if pretty logging is enabled by checking the LOG_PRETTY environment variable.
// Returns true if the environment variable is set to any non-empty value.
// Returns false when the variable is unset or set to an empty string.
//
// This provides an opt-in mechanism for human-readable formatted logging output,
// which is useful for development and debugging scenarios.
func IsPrettyLogEnabled() bool {
	v, ok := os.LookupEnv("LOG_PRETTY")
	if ok && v != "" {
		return true
	}

	return false
}
