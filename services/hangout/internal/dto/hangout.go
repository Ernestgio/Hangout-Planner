package dto

type CreateHangoutRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"datetime=2006-01-02T15:04:05Z07:00:00.000Z"`
	Status      string `json:"status" validate:"oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
}

type UpdateHangoutRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"datetime=2006-01-02T15:04:05Z07:00:00.000Z"`
	Status      string `json:"status" validate:"oneof=PLANNING CONFIRMED EXECUTED CANCELLED"`
}
