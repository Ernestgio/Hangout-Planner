package dto

import "github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/constants"

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type StandardResponseBuilder struct {
	env string
}

func NewStandardResponseBuilder(env string) *StandardResponseBuilder {
	return &StandardResponseBuilder{
		env: env,
	}
}

func (b *StandardResponseBuilder) NewSuccessResponse(message string, data interface{}) *StandardResponse {
	return &StandardResponse{
		Status:  constants.SuccessStatus,
		Message: message,
		Data:    data,
	}
}

func (b *StandardResponseBuilder) NewErrorResponse(err error) *StandardResponse {
	if b.env == constants.ProductionEnv {
		return &StandardResponse{
			Status:  constants.ErrorStatus,
			Message: "An unexpected error occurred. Please try again later.",
		}
	}

	return &StandardResponse{
		Status:  constants.ErrorStatus,
		Message: err.Error(),
	}
}
