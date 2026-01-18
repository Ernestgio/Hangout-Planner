package config

import (
	"os"
	"strconv"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	AppName    string
	AppPort    string
	DBConfig   *DBConfig
	S3Config   *S3Config
	OTELConfig *OTELConfig
	MTLSConfig *MTLSConfig
}

func Load() (*Config, error) {
	if os.Getenv("ENV") != constants.ProductionEnv {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Env:        getEnv("ENV", constants.DevEnv),
		AppName:    getEnv("APP_NAME", constants.DefaultAppName),
		AppPort:    getEnv("APP_PORT", constants.DefaultAppPort),
		DBConfig:   NewDBConfig(),
		S3Config:   NewS3Config(),
		OTELConfig: NewOTELConfig(),
		MTLSConfig: NewMTLSConfig(),
	}

	if cfg.AppPort == "" {
		return nil, apperrors.ErrAppPortRequired
	}
	if cfg.Env == constants.ProductionEnv && cfg.DBConfig.DBPassword == "" {
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

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
