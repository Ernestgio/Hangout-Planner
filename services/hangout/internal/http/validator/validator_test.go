package validator_test

import (
	"testing"

	. "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCustomValidator(t *testing.T) {
	type user struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	testCases := []struct {
		name        string
		payload     interface{}
		expectError bool
	}{
		{
			name: "success: valid payload",
			payload: &user{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expectError: false,
		},
		{
			name: "error: invalid payload with missing field",
			payload: &user{
				Email: "john.doe@example.com",
			},
			expectError: true,
		},
		{
			name: "error: invalid payload with bad format",
			payload: &user{
				Name:  "John Doe",
				Email: "not-an-email",
			},
			expectError: true,
		},
		{
			name:        "error: non-struct payload",
			payload:     "a string",
			expectError: true,
		},
	}

	t.Run("NewValidator should return a valid validator instance", func(t *testing.T) {
		v := NewValidator()
		require.NotNil(t, v)
		require.Implements(t, (*echo.Validator)(nil), v)
	})

	validator := NewValidator()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.payload)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
