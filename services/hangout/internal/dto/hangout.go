package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateHangoutRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Date        string  `json:"date" validate:"datetime=2006-01-02T15:04:05Z07:00:00.000Z"`
	Status      *string `json:"status" validate:"oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
}

type UpdateHangoutRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"datetime=2006-01-02T15:04:05Z07:00:00.000Z"`
	Status      string `json:"status" validate:"oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
}

type HangoutDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
type HangoutListItemResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Date      string    `json:"date"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
}
