package services

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
)

type AuthService interface {
	SignUser(ctx context.Context, request *dto.SignUpRequest) (*domain.User, error)
	SignInUser(ctx context.Context, request *dto.SignInRequest) (*dto.SignInResponse, error)
}

type authService struct {
	userService UserService
	jwtUtils    utils.JWTUtils
	bcrytpUtils utils.BcryptUtils
	metrics     *otel.MetricsRecorder
}

func NewAuthService(userService UserService, jwtUtils utils.JWTUtils, bcrytpUtils utils.BcryptUtils, metrics *otel.MetricsRecorder) AuthService {
	return &authService{
		userService: userService,
		jwtUtils:    jwtUtils,
		bcrytpUtils: bcrytpUtils,
		metrics:     metrics,
	}
}

func (s *authService) SignUser(ctx context.Context, request *dto.SignUpRequest) (*domain.User, error) {
	start := time.Now()

	user, err := s.userService.CreateUser(ctx, dto.CreateUserRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		s.metrics.RecordAuth(ctx, "signup", "error", time.Since(start))
		return nil, err
	}

	s.metrics.RecordAuth(ctx, "signup", "success", time.Since(start))
	return user, nil
}

func (s *authService) SignInUser(ctx context.Context, request *dto.SignInRequest) (*dto.SignInResponse, error) {
	start := time.Now()

	user, err := s.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		s.metrics.RecordAuth(ctx, "signin", "error", time.Since(start))
		return nil, err
	}
	if user == nil {
		s.metrics.RecordAuth(ctx, "signin", "error", time.Since(start))
		return nil, apperrors.ErrInvalidCredentials
	}

	err = s.bcrytpUtils.CompareHashAndPassword(user.Password, request.Password)
	if err != nil {
		s.metrics.RecordAuth(ctx, "signin", "error", time.Since(start))
		return nil, err
	}

	token, err := s.jwtUtils.Generate(user)
	if err != nil {
		s.metrics.RecordAuth(ctx, "signin", "error", time.Since(start))
		return nil, err
	}

	s.metrics.RecordAuth(ctx, "signin", "success", time.Since(start))
	return &dto.SignInResponse{
		Token: token,
	}, nil
}
