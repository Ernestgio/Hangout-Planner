package mapping

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/dto"
	models "github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/models"
)

func UserCreateRequestToModel(request dto.UserCreateRequest) models.User {
	return models.User{
		Name:  request.Name,
		Email: request.Email,
	}
}
