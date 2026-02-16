package config

import (
	"strconv"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
)

type OTELConfig struct {
	ServiceName     string
	ServiceVersion  string
	Environment     string
	Endpoint        string
	TraceEndpoint   string
	UseStdout       bool
	TraceSampleRate float64
}

func NewOTELConfig() *OTELConfig {
	traceEndpoint := getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	metricsEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", constants.DefaultOTELEndpoint)

	if traceEndpoint == "" {
		traceEndpoint = metricsEndpoint
	}

	sampleRate := constants.DefaultTraceSampleRate
	if sampleRateStr := getEnv("OTEL_TRACE_SAMPLE_RATE", ""); sampleRateStr != "" {
		if rate, err := strconv.ParseFloat(sampleRateStr, 64); err == nil {
			sampleRate = rate
		}
	}

	return &OTELConfig{
		ServiceName:     getEnv("OTEL_SERVICE_NAME", constants.DefaultAppName),
		ServiceVersion:  getEnv("OTEL_SERVICE_VERSION", constants.DefaultOTELServiceVersion),
		Environment:     getEnv("ENV", constants.DevEnv),
		Endpoint:        metricsEndpoint,
		TraceEndpoint:   traceEndpoint,
		UseStdout:       getEnv("OTEL_USE_STDOUT", constants.DefaultOTELUseStdout) == "true",
		TraceSampleRate: sampleRate,
	}
}
