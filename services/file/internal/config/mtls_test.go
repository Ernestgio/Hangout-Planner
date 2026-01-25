package config_test

import (
	"os"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewMTLSConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *config.MTLSConfig
	}{
		{
			name: "WithEnvVars_EnabledTrue",
			env: map[string]string{
				"MTLS_ENABLED":   "true",
				"MTLS_CERT_FILE": "/custom/cert.crt",
				"MTLS_KEY_FILE":  "/custom/key.key",
				"MTLS_CA_FILE":   "/custom/ca.crt",
			},
			expected: &config.MTLSConfig{
				Enabled:  true,
				CertFile: "/custom/cert.crt",
				KeyFile:  "/custom/key.key",
				CAFile:   "/custom/ca.crt",
			},
		},
		{
			name: "WithEnvVars_EnabledFalse",
			env: map[string]string{
				"MTLS_ENABLED":   "false",
				"MTLS_CERT_FILE": "/path/to/cert.crt",
				"MTLS_KEY_FILE":  "/path/to/key.key",
				"MTLS_CA_FILE":   "/path/to/ca.crt",
			},
			expected: &config.MTLSConfig{
				Enabled:  false,
				CertFile: "/path/to/cert.crt",
				KeyFile:  "/path/to/key.key",
				CAFile:   "/path/to/ca.crt",
			},
		},
		{
			name: "WithoutEnvVars_UseDefaults",
			env:  map[string]string{},
			expected: &config.MTLSConfig{
				Enabled:  false,
				CertFile: "/app/certs/file-server.crt",
				KeyFile:  "/app/certs/file-server.key",
				CAFile:   "/app/certs/ca.crt",
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

			cfg := config.NewMTLSConfig()
			require.Equal(t, tt.expected, cfg)
		})
	}
}
