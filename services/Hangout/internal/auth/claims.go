package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenCustomClaims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}
