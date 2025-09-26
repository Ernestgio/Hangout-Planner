package config

import "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"

type JwtConfig struct {
	JWTSecret          string
	JWTExpirationHours int
}

func NewJwtConfig() *JwtConfig {
	return &JwtConfig{
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", constants.DefaultJWTExpirationHours),
	}
}
