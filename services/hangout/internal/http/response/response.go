package response

import "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"

type StandardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Builder struct {
	isProduction bool
}

func NewBuilder(isProduction bool) *Builder {
	return &Builder{
		isProduction: isProduction,
	}
}

func (b *Builder) Success(message string, data any) *StandardResponse {
	return &StandardResponse{
		Status:  constants.SuccessStatus,
		Message: message,
		Data:    data,
	}
}

func (b *Builder) Error(err error) *StandardResponse {
	if b.isProduction {
		return &StandardResponse{
			Status:  constants.ErrorStatus,
			Message: constants.ProdErrorMessage,
		}
	}

	return &StandardResponse{
		Status:  constants.ErrorStatus,
		Message: err.Error(),
	}
}
