package mapper

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
)

func CreateUserRequestToModel(request dto.CreateUserRequest) domain.User {
	return domain.User{
		Name:  request.Name,
		Email: request.Email,
	}
}

func UserToResponseDTO(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
