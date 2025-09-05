package config

import (
	"errors"
	"os"

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
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Env:        getEnv("ENV", "dev"),
		AppName:    getEnv("APP_NAME", "Hangout"),
		AppPort:    getEnv("APP_PORT", "9000"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "hangout"),
	}

	if cfg.AppPort == "" {
		return nil, errors.New("APP_PORT is required")
	}
	if cfg.Env == "production" && cfg.DBPassword == "" {
		return nil, errors.New("DB_PASSWORD required in production")
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
