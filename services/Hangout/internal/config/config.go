package config

import (
	"os"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/constants"
	"github.com/joho/godotenv"
)

type Config struct {
	Env               string
	AppName           string
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	MySQLRootPassword string
}

func Load() (*Config, error) {
	if os.Getenv("ENV") != constants.ProductionEnv {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Env:               getEnv("ENV", constants.DevEnv),
		AppName:           getEnv("APP_NAME", constants.DefaultAppName),
		AppPort:           getEnv("APP_PORT", constants.DefaultAppPort),
		DBHost:            getEnv("DB_HOST", constants.DefaultDBHost),
		DBPort:            getEnv("DB_PORT", constants.DefaultDBPort),
		DBUser:            getEnv("DB_USER", constants.DefaultDBUser),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", constants.DefaultDBName),
		MySQLRootPassword: getEnv("MYSQL_ROOT_PASSWORD", ""),
	}

	if cfg.AppPort == "" {
		return nil, apperrors.ErrAppPortRequired
	}
	if cfg.Env == constants.ProductionEnv && cfg.DBPassword == "" {
		return nil, apperrors.ErrDbPasswordRequired
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
