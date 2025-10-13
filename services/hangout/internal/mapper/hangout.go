package mapper

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
)

func HangoutCreateRequestToModel(request *dto.CreateHangoutRequest) (*domain.Hangout, error) {
	parsedDate, err := time.Parse(constants.DateFormat, request.Date)
	if err != nil {
		return nil, err
	}

	return &domain.Hangout{
		Title:       request.Title,
		Description: request.Description,
		Date:        parsedDate,
		Status:      request.Status,
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
	return &dto.HangoutDetailResponse{
		ID:          hangout.ID,
		Title:       hangout.Title,
		Description: hangout.Description,
		Date:        hangout.Date,
		Status:      hangout.Status,
		CreatedAt:   hangout.CreatedAt,
	}
}

func HangoutToListItemResponseDTO(hangout *domain.Hangout) *dto.HangoutListItemResponse {
	return &dto.HangoutListItemResponse{
		ID:        hangout.ID,
		Title:     hangout.Title,
		Date:      hangout.Date,
		Status:    hangout.Status,
		CreatedAt: hangout.CreatedAt,
	}
}
