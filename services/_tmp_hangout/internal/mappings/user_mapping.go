package mappings

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
)

func UserCreateRequestToModel(request dto.UserCreateRequest) models.User {
	return models.User{
		Name:  request.Name,
		Email: request.Email,
	}
}

func UserToResponseDTO(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
