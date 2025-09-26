package utils_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestBcryptUtils(t *testing.T) {
	const correctPassword = "StrongPassword123"
	const wrongPassword = "WrongPassword456"
	const badHashFormat = "a"

	bcryptUtils := utils.NewBcryptUtils(12)

	bcryptUtilsInvalidCost := utils.NewBcryptUtils(32)

	tests := []struct {
		name          string
		run           func() error
		expectedError error
		expectGeneric bool
	}{
		{
			name: "GenerateFromPassword_Success",
			run: func() error {
				_, err := bcryptUtils.GenerateFromPassword(correctPassword)
				return err
			},
			expectedError: nil,
		},
		{
			name: "GenerateFromPassword_Failure_InvalidCost",
			run: func() error {
				_, err := bcryptUtilsInvalidCost.GenerateFromPassword(correctPassword)
				return err
			},
			expectedError: nil,
			expectGeneric: true,
		},
		{
			name: "CompareHashAndPassword_Success",
			run: func() error {
				hash, _ := bcryptUtils.GenerateFromPassword(correctPassword)
				return bcryptUtils.CompareHashAndPassword(hash, correctPassword)
			},
			expectedError: nil,
		},
		{
			name: "CompareHashAndPassword_Mismatch_AppError",
			run: func() error {
				hash, _ := bcryptUtils.GenerateFromPassword(correctPassword)
				return bcryptUtils.CompareHashAndPassword(hash, wrongPassword)
			},
			expectedError: apperrors.ErrInvalidCredentials,
		},
		{
			name: "CompareHashAndPassword_BadHash_GenericError",
			run: func() error {
				return bcryptUtils.CompareHashAndPassword(badHashFormat, correctPassword)
			},
			expectGeneric: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if tt.expectedError == nil && !tt.expectGeneric {
				require.NoError(t, err)
			} else if tt.expectGeneric {
				require.Error(t, err)

				if tt.name == "GenerateFromPassword_Failure_InvalidCost" {
					require.Contains(t, err.Error(), "outside allowed range", "Expected invalid cost error, but got: %v", err)
				} else {
					require.False(t, errors.Is(err, apperrors.ErrInvalidCredentials))
					require.Contains(t, err.Error(), "bcrypt:", "Expected generic bcrypt error, but got: %v", err)
				}
			} else {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			}
		})
	}
}
