package utils

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTUtils interface {
	Generate(user *domain.User) (string, error)
}

type jwtUtils struct {
	Secret          string
	ExpirationHours int
}

func NewJWTUtils(jWtConfig *config.JwtConfig) JWTUtils {
	return &jwtUtils{Secret: jWtConfig.JWTSecret, ExpirationHours: jWtConfig.JWTExpirationHours}
}

func (j *jwtUtils) Generate(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.ExpirationHours) * time.Hour)
	claims := &auth.TokenCustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
