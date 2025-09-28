package config_test

import (
	"os"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestNewJwtConfig(t *testing.T) {
	tests := []struct {
		name             string
		envSecret        string
		envExpiration    string
		expectedSecret   string
		expectedExpHours int
	}{
		{
			name:             "WithEnvVars",
			envSecret:        "super-secret",
			envExpiration:    "72",
			expectedSecret:   "super-secret",
			expectedExpHours: 72,
		},
		{
			name:             "WithoutEnvVars_UseDefaults",
			envSecret:        "",
			envExpiration:    "",
			expectedSecret:   "",
			expectedExpHours: constants.DefaultJWTExpirationHours,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envSecret != "" {
				if err := os.Setenv("JWT_SECRET", tt.envSecret); err != nil {
					t.Fatalf("failed to set JWT_SECRET: %v", err)
				}
			} else {
				if err := os.Unsetenv("JWT_SECRET"); err != nil {
					t.Fatalf("failed to unset JWT_SECRET: %v", err)
				}
			}

			if tt.envExpiration != "" {
				if err := os.Setenv("JWT_EXPIRATION_HOURS", tt.envExpiration); err != nil {
					t.Fatalf("failed to set JWT_EXPIRATION_HOURS: %v", err)
				}
			} else {
				if err := os.Unsetenv("JWT_EXPIRATION_HOURS"); err != nil {
					t.Fatalf("failed to unset JWT_EXPIRATION_HOURS: %v", err)
				}
			}

			cfg := config.NewJwtConfig()

			require.Equal(t, tt.expectedSecret, cfg.JWTSecret)
			require.Equal(t, tt.expectedExpHours, cfg.JWTExpirationHours)
		})
	}
}
