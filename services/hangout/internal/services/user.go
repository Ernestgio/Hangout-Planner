package services

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
)

type UserService interface {
	CreateUser(request dto.UserCreateRequest) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type userService struct {
	userRepo    repository.UserRepository
	bcryptUtils utils.BcryptUtils
}

func NewUserService(userRepo repository.UserRepository, bcryptUtils utils.BcryptUtils) UserService {
	return &userService{userRepo: userRepo, bcryptUtils: bcryptUtils}
}

func (s *userService) CreateUser(request dto.UserCreateRequest) (*domain.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(request.Email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	user := mapper.UserCreateRequestToModel(request)
	hashedPassword, err := s.bcryptUtils.GenerateFromPassword(request.Password)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(email)
}
