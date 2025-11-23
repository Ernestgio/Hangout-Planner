package mapper

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/constants"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
)

func HangoutCreateRequestToModel(request *dto.CreateHangoutRequest) (*domain.Hangout, error) {
	parsedDate, err := time.Parse(constants.DateFormat, request.Date)
	if err != nil {
		return nil, err
	}

	activities := make([]*domain.Activity, len(request.ActivityIDs))
	for i, id := range request.ActivityIDs {
		activities[i] = &domain.Activity{ID: id}
	}

	return &domain.Hangout{
		Title:       request.Title,
		Description: request.Description,
		Date:        parsedDate,
		Status:      request.Status,
		Activities:  activities,
	}, nil
}

func ApplyUpdateToHangout(hangout *domain.Hangout, req *dto.UpdateHangoutRequest) error {
	if req.Title != "" {
		hangout.Title = req.Title
	}
	hangout.Status = enums.HangoutStatus(req.Status)
	parsedDate, err := time.Parse(constants.DateFormat, req.Date)
	if err != nil {
		return err
	}
	hangout.Date = parsedDate

	if req.Description != nil {
		hangout.Description = req.Description
	}

	return nil
}

func HangoutToDetailResponseDTO(hangout *domain.Hangout) *dto.HangoutDetailResponse {
	if hangout == nil {
		return nil
	}

	activityDTOs := make([]dto.ActivityTagResponse, len(hangout.Activities))
	for i, act := range hangout.Activities {
		activityDTOs[i] = dto.ActivityTagResponse{
			ID:   act.ID,
			Name: act.Name,
		}
	}

	return &dto.HangoutDetailResponse{
		ID:          hangout.ID,
		Title:       hangout.Title,
		Description: hangout.Description,
		Date:        types.JSONTime(hangout.Date),
		Status:      hangout.Status,
		CreatedAt:   types.JSONTime(hangout.CreatedAt),
		Activities:  activityDTOs,
	}
}

func HangoutToListItemResponseDTO(hangout *domain.Hangout) *dto.HangoutListItemResponse {
	if hangout == nil {
		return nil
	}

	return &dto.HangoutListItemResponse{
		ID:        hangout.ID,
		Title:     hangout.Title,
		Date:      types.JSONTime(hangout.Date),
		Status:    hangout.Status,
		CreatedAt: types.JSONTime(hangout.CreatedAt),
	}
}

func HangoutsToListItemResponseDTOs(hangouts []domain.Hangout) []*dto.HangoutListItemResponse {
	if hangouts == nil {
		return make([]*dto.HangoutListItemResponse, 0)
	}
	responses := make([]*dto.HangoutListItemResponse, len(hangouts))
	for i, hangout := range hangouts {
		responses[i] = HangoutToListItemResponseDTO(&hangout)
	}
	return responses
}
