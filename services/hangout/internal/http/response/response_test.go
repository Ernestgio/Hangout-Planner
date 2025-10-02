package response_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	. "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/stretchr/testify/require"
)

func TestBuilderSuccess(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}

	testCases := []struct {
		name            string
		message         string
		data            any
		expectedReponse *StandardResponse
	}{
		{
			name:    "with struct data",
			message: "User found",
			data:    user{ID: 1, Name: "John Doe"},
			expectedReponse: &StandardResponse{
				Status:  constants.SuccessStatus,
				Message: "User found",
				Data:    user{ID: 1, Name: "John Doe"},
			},
		},
		{
			name:    "with nil data",
			message: "Resource deleted",
			data:    nil,
			expectedReponse: &StandardResponse{
				Status:  constants.SuccessStatus,
				Message: "Resource deleted",
				Data:    nil,
			},
		},
	}

	builder := NewBuilder(false)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResponse := builder.Success(tc.message, tc.data)
			require.Equal(t, tc.expectedReponse, actualResponse)
		})
	}
}

func TestBuilderError(t *testing.T) {
	specificError := errors.New("a specific database error occurred")

	t.Run("in production environment", func(t *testing.T) {
		builder := NewBuilder(true)
		expectedResponse := &StandardResponse{
			Status:  constants.ErrorStatus,
			Message: constants.ProdErrorMessage,
		}

		actualResponse := builder.Error(specificError)
		require.Equal(t, expectedResponse.Status, actualResponse.Status)
		require.Equal(t, expectedResponse.Message, actualResponse.Message)
		require.Nil(t, actualResponse.Data)
	})

	t.Run("in development environment", func(t *testing.T) {
		builder := NewBuilder(false)
		actualResponse := builder.Error(specificError)

		require.Equal(t, constants.ErrorStatus, actualResponse.Status)
		require.NotEmpty(t, actualResponse.Message)
		require.NotEqual(t, constants.ProdErrorMessage, actualResponse.Message)
		require.Nil(t, actualResponse.Data)
	})
}
