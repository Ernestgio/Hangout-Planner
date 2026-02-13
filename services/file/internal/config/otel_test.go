package config_test

import (
	"os"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestNewOTELConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *config.OTELConfig
	}{
		{
			name: "WithEnvVars_EnabledTrue_StdoutTrue",
			env: map[string]string{
				"OTEL_ENABLED":                "true",
				"OTEL_EXPORTER_OTLP_ENDPOINT": "otel.example.com:4318",
				"OTEL_USE_STDOUT":             "true",
				"SERVICE_VERSION":             "2.0.0",
			},
			expected: &config.OTELConfig{
				Enabled:         true,
				Endpoint:        "otel.example.com:4318",
				TraceEndpoint:   "otel.example.com:4318",
				UseStdout:       true,
				ServiceVersion:  "2.0.0",
				TraceSampleRate: constants.DefaultTraceSampleRate,
			},
		},
		{
			name: "WithEnvVars_EnabledFalse_StdoutFalse",
			env: map[string]string{
				"OTEL_ENABLED":                "false",
				"OTEL_EXPORTER_OTLP_ENDPOINT": "localhost:9090",
				"OTEL_USE_STDOUT":             "false",
				"SERVICE_VERSION":             "3.5.1",
			},
			expected: &config.OTELConfig{
				Enabled:         false,
				Endpoint:        "localhost:9090",
				TraceEndpoint:   "localhost:9090",
				UseStdout:       false,
				ServiceVersion:  "3.5.1",
				TraceSampleRate: constants.DefaultTraceSampleRate,
			},
		},
		{
			name: "WithoutEnvVars_UseDefaults",
			env:  map[string]string{},
			expected: &config.OTELConfig{
				Enabled:         true,
				Endpoint:        constants.DefaultOTELEndpoint,
				TraceEndpoint:   constants.DefaultOTELEndpoint,
				UseStdout:       false,
				ServiceVersion:  "1.0.0",
				TraceSampleRate: constants.DefaultTraceSampleRate,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env var %s: %v", k, err)
				}
			}

			cfg := config.NewOTELConfig()
			require.Equal(t, tt.expected, cfg)
		})
	}
}
