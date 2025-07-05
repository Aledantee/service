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
