package config_test

import (
	"os"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestNewDBConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *config.DBConfig
	}{
		{
			name: "WithEnvVars",
			env: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "admin",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "mydb",
			},
			expected: &config.DBConfig{
				DBHost:     "db.example.com",
				DBPort:     "5433",
				DBUser:     "admin",
				DBPassword: "secret",
				DBName:     "mydb",
			},
		},
		{
			name: "WithoutEnvVars_UseDefaults",
			env:  map[string]string{},
			expected: &config.DBConfig{
				DBHost:     constants.DefaultDBHost,
				DBPort:     constants.DefaultDBPort,
				DBUser:     constants.DefaultDBUser,
				DBPassword: "",
				DBName:     constants.DefaultDBName,
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

			cfg := config.NewDBConfig()
			require.Equal(t, tt.expected, cfg)
		})
	}
}
