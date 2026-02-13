package config

import (
	"strconv"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
)

type OTELConfig struct {
	Enabled         bool
	Endpoint        string
	TraceEndpoint   string
	UseStdout       bool
	ServiceVersion  string
	TraceSampleRate float64
}

func NewOTELConfig() *OTELConfig {
	enabled := getEnv("OTEL_ENABLED", constants.DefaultOtelEnabled) == "true"
	useStdout := getEnv("OTEL_USE_STDOUT", constants.DefaultOtelUseStdout) == "true"

	endpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", constants.DefaultOTELEndpoint)
	traceEndpoint := getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", endpoint)

	sampleRateStr := getEnv("OTEL_TRACE_SAMPLE_RATE", strconv.FormatFloat(constants.DefaultTraceSampleRate, 'f', 1, 64))
	sampleRate, err := strconv.ParseFloat(sampleRateStr, 64)
	if err != nil || sampleRate < 0 || sampleRate > 1 {
		sampleRate = constants.DefaultTraceSampleRate
	}

	return &OTELConfig{
		Enabled:         enabled,
		Endpoint:        endpoint,
		TraceEndpoint:   traceEndpoint,
		UseStdout:       useStdout,
		ServiceVersion:  getEnv("SERVICE_VERSION", constants.DefaultServiceVersion),
		TraceSampleRate: sampleRate,
	}
}
