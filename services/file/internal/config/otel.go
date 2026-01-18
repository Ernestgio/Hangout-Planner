package config

import "github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"

type OTELConfig struct {
	Enabled        bool
	Endpoint       string
	UseStdout      bool
	ServiceVersion string
}

func NewOTELConfig() *OTELConfig {
	enabled := getEnv("OTEL_ENABLED", "true") == "true"
	useStdout := getEnv("OTEL_USE_STDOUT", "false") == "true"

	return &OTELConfig{
		Enabled:        enabled,
		Endpoint:       getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", constants.DefaultOTELEndpoint),
		UseStdout:      useStdout,
		ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
	}
}
