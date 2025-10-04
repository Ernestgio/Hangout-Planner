package mapper

import (
	"time"

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
	}, nil
}

func HangoutToResponseDTO(hangout *domain.Hangout) *dto.HangoutDetailResponse {
	return &dto.HangoutDetailResponse{
		ID:          hangout.ID,
		Title:       hangout.Title,
		Description: hangout.Description,
		Date:        hangout.Date,
		Status:      hangout.Status,
	}
}
