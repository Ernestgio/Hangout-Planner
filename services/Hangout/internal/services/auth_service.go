package services

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
)

type AuthService interface {
	SignUser(request *dto.SignUpRequest) (*models.User, error)
	SignInUser(request *dto.SignInRequest) (*dto.SignInResponse, error)
}

type authService struct {
	userService UserService
	jwtUtils    utils.JWTUtils
	bcrytpUtils utils.BcryptUtils
}

func NewAuthService(userService UserService, jwtUtils utils.JWTUtils, bcrytpUtils utils.BcryptUtils) AuthService {
	return &authService{
		userService: userService,
		jwtUtils:    jwtUtils,
		bcrytpUtils: bcrytpUtils,
	}
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

func (s *authService) SignInUser(request *dto.SignInRequest) (*dto.SignInResponse, error) {
	user, err := s.userService.GetUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	err = s.bcrytpUtils.CompareHashAndPassword(user.Password, request.Password)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtUtils.Generate(user)
	if err != nil {
		return nil, err
	}

	return &dto.SignInResponse{
		Token: token,
	}, nil
}
