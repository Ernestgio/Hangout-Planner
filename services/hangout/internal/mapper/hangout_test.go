package mapper_test

import (
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/constants"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
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
	activityID1 := uuid.New()
	activityID2 := uuid.New()

	testCases := []struct {
		name        string
		request     *dto.CreateHangoutRequest
		expectError bool
		checkResult func(t *testing.T, hangout *domain.Hangout, err error)
	}{
		{
			name: "success with activities",
			request: &dto.CreateHangoutRequest{
				Title:       "Test Hangout",
				Description: stringPtr("A cool event."),
				Date:        validTimeStr,
				Status:      enums.StatusPlanning,
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
			},
			expectError: false,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, hangout)
				require.Equal(t, "Test Hangout", hangout.Title)
				require.Equal(t, "A cool event.", *hangout.Description)
				require.Equal(t, parsedTime, hangout.Date)
				require.Equal(t, enums.StatusPlanning, hangout.Status)
				require.Len(t, hangout.Activities, 2)
				require.Equal(t, activityID1, hangout.Activities[0].ID)
				require.Equal(t, activityID2, hangout.Activities[1].ID)
			},
		},
		{
			name: "success without activities",
			request: &dto.CreateHangoutRequest{
				Title:       "Test Hangout No Activities",
				Description: stringPtr(""),
				Date:        validTimeStr,
				Status:      enums.StatusConfirmed,
				ActivityIDs: []uuid.UUID{},
			},
			expectError: false,
			checkResult: func(t *testing.T, hangout *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, hangout)
				require.Empty(t, hangout.Activities)
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

func TestApplyUpdateToHangout(t *testing.T) {
	initialDate := time.Now().Add(-24 * time.Hour)
	newDateStr := "2025-12-25 18:00:00.000"
	parsedNewDate, _ := time.Parse(constants.DateFormat, newDateStr)

	testCases := []struct {
		name           string
		initialHangout *domain.Hangout
		request        *dto.UpdateHangoutRequest
		expectError    bool
		checkResult    func(t *testing.T, hangout *domain.Hangout)
	}{
		{
			name: "success: full update of basic fields",
			initialHangout: &domain.Hangout{
				ID:    uuid.New(),
				Title: "Old Title",
				Date:  initialDate,
			},
			request: &dto.UpdateHangoutRequest{
				Title:       "New Title",
				Description: stringPtr("New Description"),
				Date:        newDateStr,
				Status:      enums.StatusConfirmed,
			},
			expectError: false,
			checkResult: func(t *testing.T, hangout *domain.Hangout) {
				require.Equal(t, "New Title", hangout.Title)
				require.Equal(t, "New Description", *hangout.Description)
				require.Equal(t, parsedNewDate, hangout.Date)
				require.Equal(t, enums.StatusConfirmed, hangout.Status)
			},
		},
		{
			name: "error: invalid date format",
			initialHangout: &domain.Hangout{
				ID: uuid.New(),
			},
			request: &dto.UpdateHangoutRequest{
				Date: "invalid-date-format",
			},
			expectError: true,
			checkResult: func(t *testing.T, hangout *domain.Hangout) {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mapper.ApplyUpdateToHangout(tc.initialHangout, tc.request)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.checkResult(t, tc.initialHangout)
			}
		})
	}
}

func TestHangoutToDetailResponseDTO(t *testing.T) {
	hangoutID := uuid.New()
	activityID1 := uuid.New()
	activityID2 := uuid.New()
	now := time.Now()

	hangoutWithActivities := &domain.Hangout{
		ID:          hangoutID,
		Title:       "Detail View",
		Description: stringPtr("Detailed description."),
		Date:        now,
		Status:      enums.StatusExecuted,
		CreatedAt:   now,
		Activities: []*domain.Activity{
			{ID: activityID1, Name: "Hiking"},
			{ID: activityID2, Name: "Coffee"},
		},
	}

	testCases := []struct {
		name        string
		input       *domain.Hangout
		checkResult func(t *testing.T, res *dto.HangoutDetailResponse)
	}{
		{
			name:  "success with activities",
			input: hangoutWithActivities,
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse) {
				require.NotNil(t, res)
				require.Equal(t, hangoutID, res.ID)
				require.Equal(t, "Detail View", res.Title)
				require.NotNil(t, res.Description)
				require.Equal(t, "Detailed description.", *res.Description)
				require.Equal(t, types.JSONTime(now), res.Date)
				require.Equal(t, enums.StatusExecuted, res.Status)
				require.Len(t, res.Activities, 2)
				require.Equal(t, activityID1, res.Activities[0].ID)
				require.Equal(t, "Hiking", res.Activities[0].Name)
				require.Equal(t, activityID2, res.Activities[1].ID)
				require.Equal(t, "Coffee", res.Activities[1].Name)
			},
		},
		{
			name: "success without activities",
			input: &domain.Hangout{
				ID:          uuid.New(),
				Title:       "No Activities",
				Date:        now,
				Description: stringPtr(""),
				Activities:  []*domain.Activity{},
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse) {
				require.NotNil(t, res)
				require.Empty(t, res.Activities)
			},
		},
		{
			name:  "nil input",
			input: nil,
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse) {
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := mapper.HangoutToDetailResponseDTO(tc.input)
			tc.checkResult(t, response)
		})
	}
}

func TestHangoutToListItemResponseDTO(t *testing.T) {
	hangoutID := uuid.New()
	now := time.Now()
	hangout := &domain.Hangout{
		ID:        hangoutID,
		Title:     "List Item View",
		Date:      now,
		Status:    enums.StatusCancelled,
		CreatedAt: now,
	}

	testCases := []struct {
		name        string
		input       *domain.Hangout
		checkResult func(t *testing.T, res *dto.HangoutListItemResponse)
	}{
		{
			name:  "success",
			input: hangout,
			checkResult: func(t *testing.T, res *dto.HangoutListItemResponse) {
				require.NotNil(t, res)
				require.Equal(t, hangoutID, res.ID)
				require.Equal(t, "List Item View", res.Title)
				require.Equal(t, types.JSONTime(now), res.Date)
				require.Equal(t, enums.StatusCancelled, res.Status)
				require.Equal(t, types.JSONTime(now), res.CreatedAt)
			},
		},
		{
			name:  "nil input",
			input: nil,
			checkResult: func(t *testing.T, res *dto.HangoutListItemResponse) {
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := mapper.HangoutToListItemResponseDTO(tc.input)
			tc.checkResult(t, response)
		})
	}
}

func TestHangoutsToListItemResponseDTOs(t *testing.T) {
	hangout1 := domain.Hangout{ID: uuid.New(), Title: "First Hangout"}
	hangout2 := domain.Hangout{ID: uuid.New(), Title: "Second Hangout"}

	testCases := []struct {
		name          string
		inputHangouts []domain.Hangout
		checkResult   func(t *testing.T, res []*dto.HangoutListItemResponse)
	}{
		{
			name:          "non-empty slice",
			inputHangouts: []domain.Hangout{hangout1, hangout2},
			checkResult: func(t *testing.T, res []*dto.HangoutListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 2)
				require.Equal(t, hangout1.ID, res[0].ID)
				require.Equal(t, hangout2.Title, res[1].Title)
			},
		},
		{
			name:          "empty slice",
			inputHangouts: []domain.Hangout{},
			checkResult: func(t *testing.T, res []*dto.HangoutListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 0)
			},
		},
		{
			name:          "nil slice",
			inputHangouts: nil,
			checkResult: func(t *testing.T, res []*dto.HangoutListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapper.HangoutsToListItemResponseDTOs(tc.inputHangouts)
			tc.checkResult(t, result)
		})
	}
}
