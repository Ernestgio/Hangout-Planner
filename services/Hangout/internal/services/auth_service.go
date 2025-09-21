package services

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
)

type AuthService interface {
	SignUser(request *dto.SignUpRequest) (*models.User, error)
}

type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{userService: userService}
}

func (s *authService) SignUser(request *dto.SignUpRequest) (*models.User, error) {
	user, err := s.userService.CreateUser(dto.UserCreateRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}
