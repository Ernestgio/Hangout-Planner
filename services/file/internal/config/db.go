package config

import "github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		DBHost:     getEnv("DB_HOST", constants.DefaultDBHost),
		DBPort:     getEnv("DB_PORT", constants.DefaultDBPort),
		DBUser:     getEnv("DB_USER", constants.DefaultDBUser),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", constants.DefaultDBName),
	}
}
