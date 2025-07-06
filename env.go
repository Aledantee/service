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

func IsDebugEnabled() bool {
	v, ok := os.LookupEnv("DEBUG")
	if ok && v != "" {
		return true
	}

	return false
}

func LogLevel() string {
	v, ok := os.LookupEnv("LOG_LEVEL")
	if ok && v != "" {
		return v
	}

	return "info"
}

func IsJsonLogEnabled() bool {
	v, ok := os.LookupEnv("LOG_JSON")
	if ok && v != "" {
		return true
	}

	return false
}

func IsPrettyLogEnabled() bool {
	v, ok := os.LookupEnv("LOG_PRETTY")
	if ok && v != "" {
		return true
	}

	return false
}
