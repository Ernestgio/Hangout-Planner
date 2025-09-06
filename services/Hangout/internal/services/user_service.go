package services

import (
	"errors"

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
	repo repository.UserRepository
}

var ErrUserAlreadyExists = errors.New("user already exists with this email")

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
func (s *userService) CreateUser(request dto.UserCreateRequest) (*models.User, error) {
	existingUser, err := s.repo.GetUserByEmail(request.Email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	user := mappings.UserCreateRequestToModel(request)
	user.ID = uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
