package utils_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

const validSecret = "test-secret-key-that-is-long-enough-for-signing-HS256"

func TestNewJWTUtils(t *testing.T) {
	secret := "test"
	expiry := 24

	jwtUtil := utils.NewJWTUtils(secret, expiry)

	if jwtUtil == nil {
		t.Fatal("NewJWTUtils returned nil, expected JWTUtils implementation")
	}
}

func TestJWTUtils_Generate(t *testing.T) {
	testUser := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tests := []struct {
		name        string
		user        *models.User
		secret      string
		expiryHours int
		expectedErr error
		checkClaims bool
	}{
		{
			name:        "Success_ValidTokenGenerated",
			user:        testUser,
			secret:      validSecret,
			expiryHours: 24,
			expectedErr: nil,
			checkClaims: true,
		},
		{
			name:        "Failure_EmptySecret",
			user:        testUser,
			secret:      "",
			expiryHours: 1,
			expectedErr: nil,
			checkClaims: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtUtil := utils.NewJWTUtils(tt.secret, tt.expiryHours)

			timeBefore := time.Now().Add(-1 * time.Second)

			tokenString, err := jwtUtil.Generate(tt.user)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Generate() error = %v, want %v", err, tt.expectedErr)
			}

			if err != nil {
				return
			}

			if tokenString == "" {
				t.Fatal("Generate() returned an empty token string on success")
			}

			if tt.checkClaims {
				claims := &auth.TokenCustomClaims{}

				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("unexpected signing method")
					}
					return []byte(tt.secret), nil
				})

				if err != nil {
					t.Fatalf("Failed to parse and verify token: %v", err)
				}
				if !token.Valid {
					t.Fatal("Generated token is not valid")
				}

				if claims.UserID != tt.user.ID {
					t.Errorf("Claim UserID got = %d, want %d", claims.UserID, tt.user.ID)
				}
				if claims.Subject != tt.user.Email {
					t.Errorf("Claim Subject got = %s, want %s", claims.Subject, tt.user.Email)
				}

				timeAfter := time.Now().Add(1 * time.Second)

				if claims.IssuedAt.Before(timeBefore) || claims.IssuedAt.After(timeAfter) {
					t.Errorf("Claim IssuedAt is outside the test window: %v", claims.IssuedAt.Time)
				}

				expectedExpirationTime := timeBefore.Add(time.Duration(tt.expiryHours) * time.Hour)
				tolerance := 5 * time.Second

				if claims.ExpiresAt.Before(expectedExpirationTime.Add(-tolerance)) || claims.ExpiresAt.After(expectedExpirationTime.Add(tolerance)) {
					t.Errorf("Claim ExpiresAt is incorrect. Expected approx %v, got %v", expectedExpirationTime, claims.ExpiresAt.Time)
				}
			}
		})
	}
}
