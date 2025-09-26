package utils

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTUtils interface {
	Generate(user *models.User) (string, error)
}

type jwtUtils struct {
	jwtSecret string
	jwtExpiry int
}

func NewJWTUtils(jwtSecret string, jwtExpiry int) JWTUtils {
	return &jwtUtils{jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

func (j *jwtUtils) Generate(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.jwtExpiry) * time.Hour)
	claims := &auth.TokenCustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(j.jwtSecret))

	return tokenString, nil
}
