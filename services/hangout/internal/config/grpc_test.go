package config_test

import (
	"os"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestNewGRPCClientConfig(t *testing.T) {
	tests := []struct {
		name              string
		envFileServiceURL string
		envMTLSEnabled    string
		envCertFile       string
		envKeyFile        string
		envCAFile         string
		expectedURL       string
		expectedMTLS      bool
		expectedCertFile  string
		expectedKeyFile   string
		expectedCAFile    string
	}{
		{
			name:              "WithEnvVars_MTLSTrue",
			envFileServiceURL: "file-service:9090",
			envMTLSEnabled:    "true",
			envCertFile:       "/custom/cert.crt",
			envKeyFile:        "/custom/key.key",
			envCAFile:         "/custom/ca.crt",
			expectedURL:       "file-service:9090",
			expectedMTLS:      true,
			expectedCertFile:  "/custom/cert.crt",
			expectedKeyFile:   "/custom/key.key",
			expectedCAFile:    "/custom/ca.crt",
		},
		{
			name:              "WithEnvVars_MTLSFalse",
			envFileServiceURL: "localhost:8080",
			envMTLSEnabled:    "false",
			envCertFile:       "/path/to/cert.crt",
			envKeyFile:        "/path/to/key.key",
			envCAFile:         "/path/to/ca.crt",
			expectedURL:       "localhost:8080",
			expectedMTLS:      false,
			expectedCertFile:  "/path/to/cert.crt",
			expectedKeyFile:   "/path/to/key.key",
			expectedCAFile:    "/path/to/ca.crt",
		},
		{
			name:              "WithoutEnvVars_UseDefaults",
			envFileServiceURL: "",
			envMTLSEnabled:    "",
			envCertFile:       "",
			envKeyFile:        "",
			envCAFile:         "",
			expectedURL:       constants.DefaultFileServiceURL,
			expectedMTLS:      true,
			expectedCertFile:  constants.DefaultMTLSCertPath,
			expectedKeyFile:   constants.DefaultMTLSKeyPath,
			expectedCAFile:    constants.DefaultMTLSCAPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envFileServiceURL != "" {
				if err := os.Setenv("FILE_SERVICE_URL", tt.envFileServiceURL); err != nil {
					t.Fatalf("failed to set FILE_SERVICE_URL: %v", err)
				}
			} else {
				if err := os.Unsetenv("FILE_SERVICE_URL"); err != nil {
					t.Fatalf("failed to unset FILE_SERVICE_URL: %v", err)
				}
			}

			if tt.envMTLSEnabled != "" {
				if err := os.Setenv("GRPC_MTLS_ENABLED", tt.envMTLSEnabled); err != nil {
					t.Fatalf("failed to set GRPC_MTLS_ENABLED: %v", err)
				}
			} else {
				if err := os.Unsetenv("GRPC_MTLS_ENABLED"); err != nil {
					t.Fatalf("failed to unset GRPC_MTLS_ENABLED: %v", err)
				}
			}

			if tt.envCertFile != "" {
				if err := os.Setenv("GRPC_MTLS_CERT_FILE", tt.envCertFile); err != nil {
					t.Fatalf("failed to set GRPC_MTLS_CERT_FILE: %v", err)
				}
			} else {
				if err := os.Unsetenv("GRPC_MTLS_CERT_FILE"); err != nil {
					t.Fatalf("failed to unset GRPC_MTLS_CERT_FILE: %v", err)
				}
			}

			if tt.envKeyFile != "" {
				if err := os.Setenv("GRPC_MTLS_KEY_FILE", tt.envKeyFile); err != nil {
					t.Fatalf("failed to set GRPC_MTLS_KEY_FILE: %v", err)
				}
			} else {
				if err := os.Unsetenv("GRPC_MTLS_KEY_FILE"); err != nil {
					t.Fatalf("failed to unset GRPC_MTLS_KEY_FILE: %v", err)
				}
			}

			if tt.envCAFile != "" {
				if err := os.Setenv("GRPC_MTLS_CA_FILE", tt.envCAFile); err != nil {
					t.Fatalf("failed to set GRPC_MTLS_CA_FILE: %v", err)
				}
			} else {
				if err := os.Unsetenv("GRPC_MTLS_CA_FILE"); err != nil {
					t.Fatalf("failed to unset GRPC_MTLS_CA_FILE: %v", err)
				}
			}

			cfg := config.NewGRPCClientConfig()

			require.Equal(t, tt.expectedURL, cfg.FileServiceURL)
			require.Equal(t, tt.expectedMTLS, cfg.MTLSEnabled)
			require.Equal(t, tt.expectedCertFile, cfg.CertFile)
			require.Equal(t, tt.expectedKeyFile, cfg.KeyFile)
			require.Equal(t, tt.expectedCAFile, cfg.CAFile)
		})
	}
}
