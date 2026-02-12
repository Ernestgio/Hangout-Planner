package services

import (
	"context"
	"errors"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
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
	metrics     *otel.MetricsRecorder
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, bcryptUtils utils.BcryptUtils, metrics *otel.MetricsRecorder) UserService {
	return &userService{
		db:          db,
		userRepo:    userRepo,
		bcryptUtils: bcryptUtils,
		metrics:     metrics,
	}
}

func (s *userService) CreateUser(ctx context.Context, request dto.CreateUserRequest) (*domain.User, error) {
	start := time.Now()
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
		s.metrics.RecordDBOperation(ctx, "transaction", "users", time.Since(start), 0)
		return nil, err
	}

	s.metrics.RecordDBOperation(ctx, "insert", "users", time.Since(start), 1)
	return createdUser, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	start := time.Now()
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	s.metrics.RecordDBOperation(ctx, "select", "users", time.Since(start), 0)
	return user, err
}
