package services

import (
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/mappings"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(request dto.UserCreateRequest) (*models.User, error)
}

type userService struct {
	userRepo   repository.UserRepository
	bcryptCost int
}

func NewUserService(userRepo repository.UserRepository, bcryptCost int) UserService {
	return &userService{userRepo: userRepo, bcryptCost: bcryptCost}
}

func (s *userService) CreateUser(request dto.UserCreateRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(request.Email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	user := mappings.UserCreateRequestToModel(request)
	user.ID = uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
