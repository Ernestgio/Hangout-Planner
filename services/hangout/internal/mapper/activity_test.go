package mapper_test

import (
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestActivitytoDetailResponseDTO(t *testing.T) {
	activityID := uuid.New()
	now := time.Now()
	activity := &domain.Activity{
		ID:        activityID,
		Name:      "Hiking",
		CreatedAt: now,
	}

	testCases := []struct {
		name        string
		input       *domain.Activity
		count       int64
		checkResult func(t *testing.T, res *dto.ActivityDetailResponse)
	}{
		{
			name:  "success",
			input: activity,
			count: 5,
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse) {
				require.NotNil(t, res)
				require.Equal(t, activityID, res.ID)
				require.Equal(t, "Hiking", res.Name)
				require.Equal(t, int64(5), res.HangoutCount)
				require.Equal(t, types.JSONTime(now), res.CreatedAt)
			},
		},
		{
			name:  "nil input",
			input: nil,
			count: 0,
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse) {
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := mapper.ActivitytoDetailResponseDTO(tc.input, tc.count)
			tc.checkResult(t, response)
		})
	}
}

func TestActivityToListItemResponseDTO(t *testing.T) {
	activity1 := repository.ActivityWithCount{
		Activity:     domain.Activity{ID: uuid.New(), Name: "Hiking"},
		HangoutCount: 3,
	}
	activity2 := repository.ActivityWithCount{
		Activity:     domain.Activity{ID: uuid.New(), Name: "Movies"},
		HangoutCount: 10,
	}

	testCases := []struct {
		name        string
		input       []repository.ActivityWithCount
		checkResult func(t *testing.T, res []dto.ActivityListItemResponse)
	}{
		{
			name:  "non-empty slice",
			input: []repository.ActivityWithCount{activity1, activity2},
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 2)
				require.Equal(t, activity1.ID, res[0].ID)
				require.Equal(t, activity1.Name, res[0].Name)
				require.Equal(t, activity1.HangoutCount, res[0].HangoutCount)
				require.Equal(t, activity2.ID, res[1].ID)
				require.Equal(t, activity2.Name, res[1].Name)
				require.Equal(t, activity2.HangoutCount, res[1].HangoutCount)
			},
		},
		{
			name:  "empty slice",
			input: []repository.ActivityWithCount{},
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 0)
			},
		},
		{
			name:  "nil slice",
			input: nil,
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse) {
				require.NotNil(t, res)
				require.Len(t, res, 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapper.ActivityToListItemResponseDTO(tc.input)
			tc.checkResult(t, result)
		})
	}
}

func TestApplyUpdateToActivity(t *testing.T) {
	testCases := []struct {
		name           string
		initialHangout *domain.Activity
		request        *dto.UpdateActivityRequest
		checkResult    func(t *testing.T, hangout *domain.Activity)
	}{
		{
			name: "success: update name",
			initialHangout: &domain.Activity{
				Name: "Old Name",
			},
			request: &dto.UpdateActivityRequest{
				Name: "New Name",
			},
			checkResult: func(t *testing.T, hangout *domain.Activity) {
				require.Equal(t, "New Name", hangout.Name)
			},
		},
		{
			name: "no update: empty name string",
			initialHangout: &domain.Activity{
				Name: "Old Name",
			},
			request: &dto.UpdateActivityRequest{
				Name: "",
			},
			checkResult: func(t *testing.T, hangout *domain.Activity) {
				require.Equal(t, "Old Name", hangout.Name)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapper.ApplyUpdateToActivity(tc.initialHangout, tc.request)
			tc.checkResult(t, tc.initialHangout)
		})
	}
}
