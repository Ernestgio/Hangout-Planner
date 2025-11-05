package mapper

import (
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
)

func ActivitytoDetailResponseDTO(activity *domain.Activity, hangoutCount int64) *dto.ActivityDetailResponse {
	if activity == nil {
		return nil
	}
	return &dto.ActivityDetailResponse{
		ID:           activity.ID,
		Name:         activity.Name,
		HangoutCount: hangoutCount,
		CreatedAt:    types.JSONTime(activity.CreatedAt),
	}
}
func ActivityToListItemResponseDTO(activities []repository.ActivityWithCount) []dto.ActivityListItemResponse {
	if activities == nil {
		return make([]dto.ActivityListItemResponse, 0)
	}

	activityList := make([]dto.ActivityListItemResponse, len(activities))
	for i, activityWithCount := range activities {
		activityList[i] = dto.ActivityListItemResponse{
			ID:           activityWithCount.ID,
			Name:         activityWithCount.Name,
			HangoutCount: activityWithCount.HangoutCount,
		}
	}

	return activityList
}

func ApplyUpdateToActivity(activity *domain.Activity, req *dto.UpdateActivityRequest) {
	if req.Name != "" {
		activity.Name = req.Name
	}
}
