package utils_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Helper to generate a valid bcrypt hash
func generateHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

// TestCompareHashAndPassword covers all three distinct return paths for 100% logic coverage
func TestCompareHashAndPassword(t *testing.T) {
	const correctPassword = "StrongPassword123"
	const wrongPassword = "WrongPassword456"

	knownGoodHash := generateHash(t, correctPassword)

	// A very short string reliably triggers a generic bcrypt format error (like "hashedSecret too short")
	const badHashFormat = "a"

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectedError  error
		expectGeneric  bool
	}{
		{
			name:           "Success_MatchingPassword",
			hashedPassword: knownGoodHash,
			password:       correctPassword,
			expectedError:  nil,
			expectGeneric:  false,
		},
		{
			name:           "Failure_MismatchedPassword_AppError",
			hashedPassword: knownGoodHash,
			password:       wrongPassword,
			expectedError:  apperrors.ErrInvalidCredentials,
			expectGeneric:  false,
		},
		{
			// This test case ensures the final 'return err' line for non-mismatch errors is hit.
			name:           "Failure_BadHashFormat_GenericBcryptError",
			hashedPassword: badHashFormat,
			password:       correctPassword,
			expectedError:  nil, // We expect *an* error, but check its type/content specifically below
			expectGeneric:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.CompareHashAndPassword(tt.hashedPassword, tt.password)

			if tt.expectedError == nil && !tt.expectGeneric {
				require.NoError(t, err)
			} else if tt.expectGeneric {
				// Path 3: Generic Bcrypt Error -> returns the raw error
				require.Error(t, err)
				require.False(t, errors.Is(err, apperrors.ErrInvalidCredentials), "Should return the raw bcrypt error, not the app error")
				require.Contains(t, err.Error(), "bcrypt:", "Expected a generic bcrypt error message")
			} else {
				// Path 2: Mismatch Error -> returns apperrors.ErrInvalidCredentials
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "Expected app error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
