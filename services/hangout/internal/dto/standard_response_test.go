package dto_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/stretchr/testify/require"
)

func TestNewSuccessResponse(t *testing.T) {
	builder := dto.NewStandardResponseBuilder(constants.DevEnv)
	message, data := "ok", "sample data"
	resp := builder.NewSuccessResponse(message, data)

	require.Equal(t, constants.SuccessStatus, resp.Status)
	require.Equal(t, message, resp.Message)
	require.Equal(t, data, resp.Data)
}

func TestNewErrorResponse(t *testing.T) {
	tests := map[string]struct {
		env      string
		err      error
		expected string
	}{
		"dev env returns raw error": {
			env:      constants.DevEnv,
			err:      errors.New("boom"),
			expected: "boom",
		},
		"prod env hides details": {
			env:      constants.ProductionEnv,
			err:      errors.New("boom"),
			expected: constants.ProdErrorMessage,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			builder := dto.NewStandardResponseBuilder(tt.env)

			resp := builder.NewErrorResponse(tt.err)

			require.Equal(t, constants.ErrorStatus, resp.Status)
			require.Equal(t, tt.expected, resp.Message)
		})
	}
}
