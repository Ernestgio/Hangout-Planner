package mapper_test

import (
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func TestHangoutCreateRequestToModel(t *testing.T) {
	validTimeStr := "2025-10-05 15:00:00.000"
	parsedTime, _ := time.Parse(constants.DateFormat, validTimeStr)

	testCases := []struct {
		name        string
		request     *dto.CreateHangoutRequest
		expectError bool
		checkResult func(t *testing.T, hangout *domain.Hangout, err error)
	}{
		{
			name: "success",
			request: &dto.CreateHangoutRequest{
				Title:       "Test Hangout",
				Description: stringPtr("A cool event."),
				Date:        validTimeStr,
				Status:      enums.StatusPlanning,
			},
			expectError: false,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, hangout)
				require.Equal(t, "Test Hangout", hangout.Title)
				require.Equal(t, "A cool event.", *hangout.Description)
				require.Equal(t, parsedTime, hangout.Date)
				require.Equal(t, enums.StatusPlanning, hangout.Status)
			},
		},
		{
			name: "invalid date format",
			request: &dto.CreateHangoutRequest{
				Date: "invalid-date",
			},
			expectError: true,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.Error(t, err)
				require.Nil(t, hangout)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hangout, err := mapper.HangoutCreateRequestToModel(tc.request)
			tc.checkResult(t, hangout, err)
		})
	}
}

func TestHangoutUpdateRequestToModel(t *testing.T) {
	validTimeStr := "2025-11-20 18:30:00.000"
	parsedTime, _ := time.Parse(constants.DateFormat, validTimeStr)

	testCases := []struct {
		name        string
		request     *dto.UpdateHangoutRequest
		expectError bool
		checkResult func(t *testing.T, hangout *domain.Hangout, err error)
	}{
		{
			name: "success",
			request: &dto.UpdateHangoutRequest{
				Title:       "Updated Hangout",
				Description: stringPtr("Updated description."),
				Date:        validTimeStr,
				Status:      enums.StatusConfirmed,
			},
			expectError: false,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, hangout)
				require.Equal(t, "Updated Hangout", hangout.Title)
				require.Equal(t, "Updated description.", *hangout.Description)
				require.Equal(t, parsedTime, hangout.Date)
				require.Equal(t, enums.StatusConfirmed, hangout.Status)
			},
		},
		{
			name: "invalid date format",
			request: &dto.UpdateHangoutRequest{
				Date: "2025/11/20",
			},
			expectError: true,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.Error(t, err)
				require.Nil(t, hangout)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hangout, err := mapper.HangoutUpdateRequestToModel(tc.request)
			tc.checkResult(t, hangout, err)
		})
	}
}

func TestHangoutToDetailResponseDTO(t *testing.T) {
	hangoutID := uuid.New()
	now := time.Now()
	hangout := &domain.Hangout{
		ID:          hangoutID,
		Title:       "Detail View",
		Description: stringPtr("This is a detailed description."),
		Date:        now,
		Status:      enums.StatusExecuted,
		CreatedAt:   now,
	}

	response := mapper.HangoutToDetailResponseDTO(hangout)

	require.NotNil(t, response)
	require.Equal(t, hangoutID, response.ID)
	require.Equal(t, "Detail View", response.Title)
	require.Equal(t, "This is a detailed description.", *response.Description) // Dereference pointer for comparison
	require.Equal(t, now, response.Date)
	require.Equal(t, enums.StatusExecuted, response.Status)
	require.Equal(t, now, response.CreatedAt)
}

func TestHangoutToListItemResponseDTO(t *testing.T) {
	hangoutID := uuid.New()
	now := time.Now()
	createdAt := now.Add(-1 * time.Hour)
	hangout := &domain.Hangout{
		ID:        hangoutID,
		Title:     "List Item View",
		Date:      now,
		Status:    enums.StatusCancelled,
		CreatedAt: createdAt,
	}

	response := mapper.HangoutToListItemResponseDTO(hangout)

	require.NotNil(t, response)
	require.Equal(t, hangoutID, response.ID)
	require.Equal(t, "List Item View", response.Title)
	require.Equal(t, now, response.Date)
	require.Equal(t, enums.StatusCancelled, response.Status)
	require.Equal(t, createdAt, response.CreatedAt)
}
