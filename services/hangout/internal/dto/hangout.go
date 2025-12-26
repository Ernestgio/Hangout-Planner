package dto

import (
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/google/uuid"
)

type CreateHangoutRequest struct {
	Title       string              `json:"title" validate:"required"`
	Description *string             `json:"description"`
	Date        string              `json:"date" validate:"required,datetime=2006-01-02 15:04:05.000"`
	Status      enums.HangoutStatus `json:"status" validate:"oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
	ActivityIDs []uuid.UUID         `json:"activity_ids" validate:"dive,uuid"`
}

type UpdateHangoutRequest struct {
	Title       string              `json:"title" validate:"required"`
	Description *string             `json:"description"`
	Date        string              `json:"date" validate:"required,datetime=2006-01-02 15:04:05.000"`
	Status      enums.HangoutStatus `json:"status" validate:"required,oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
	ActivityIDs []uuid.UUID         `json:"activities" validate:"dive,uuid"`
}

type HangoutDetailResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description *string               `json:"description"`
	Date        types.JSONTime        `json:"date"`
	Status      enums.HangoutStatus   `json:"status"`
	CreatedAt   types.JSONTime        `json:"created_at"`
	Activities  []ActivityTagResponse `json:"activities"`
}

type HangoutListItemResponse struct {
	ID        uuid.UUID           `json:"id"`
	Title     string              `json:"title"`
	Date      types.JSONTime      `json:"date"`
	Status    enums.HangoutStatus `json:"status"`
	CreatedAt types.JSONTime      `json:"created_at"`
}

type PaginatedHangouts struct {
	Data       []*HangoutListItemResponse `json:"data"`
	NextCursor *uuid.UUID                 `json:"next_cursor"`
	HasMore    bool                       `json:"has_more"`
}
