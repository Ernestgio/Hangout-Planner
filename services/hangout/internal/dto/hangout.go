package dto

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/google/uuid"
)

type CreateHangoutRequest struct {
	Title       string              `json:"title" validate:"required"`
	Description *string             `json:"description"`
	Date        string              `json:"date" validate:"datetime=2006-01-02 15:04:05.000"`
	Status      enums.HangoutStatus `json:"status"`
}

type UpdateHangoutRequest struct {
	Title       string              `json:"title" validate:"required"`
	Description *string             `json:"description"`
	Date        string              `json:"date" validate:"datetime=2006-01-02 15:04:05.000"`
	Status      enums.HangoutStatus `json:"status"`
}

type HangoutDetailResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Description *string             `json:"description"`
	Date        time.Time           `json:"date"`
	Status      enums.HangoutStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
}
type HangoutListItemResponse struct {
	ID        uuid.UUID           `json:"id"`
	Title     string              `json:"title"`
	Date      time.Time           `json:"date"`
	Status    enums.HangoutStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}
