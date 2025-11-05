package dto

import (
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/google/uuid"
)

type CreateActivityRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateActivityRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type ActivityListItemResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	HangoutCount int64     `json:"hangout_count"`
}

type ActivityDetailResponse struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	HangoutCount int64          `json:"hangout_count"`
	CreatedAt    types.JSONTime `json:"created_at"`
}
