package services

import (
	"context"
	"errors"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userService struct {
	db          *gorm.DB
	userRepo    repository.UserRepository
	bcryptUtils utils.BcryptUtils
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, bcryptUtils utils.BcryptUtils) UserService {
	return &userService{db: db, userRepo: userRepo, bcryptUtils: bcryptUtils}
}

func (s *userService) CreateUser(ctx context.Context, request dto.CreateUserRequest) (*domain.User, error) {
	var createdUser *domain.User

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.userRepo.WithTx(tx)

		_, err := txRepo.GetUserByEmail(ctx, request.Email)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserAlreadyExists
		}

		user := mapper.CreateUserRequestToModel(request)
		hashedPassword, err := s.bcryptUtils.GenerateFromPassword(request.Password)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		if err := txRepo.CreateUser(ctx, &user); err != nil {
			return err
		}

		createdUser = &user
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}
