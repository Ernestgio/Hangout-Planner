package config

import "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"

type OTELConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string
	UseStdout      bool
}

func NewOTELConfig() *OTELConfig {
	return &OTELConfig{
		ServiceName:    getEnv("OTEL_SERVICE_NAME", constants.DefaultAppName),
		ServiceVersion: getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
		Environment:    getEnv("ENV", constants.DevEnv),
		Endpoint:       getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", constants.DefaultOTELEndpoint),
		UseStdout:      getEnv("OTEL_USE_STDOUT", "false") == "true",
	}
}
