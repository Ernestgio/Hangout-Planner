package services_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserService_CreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBcrypt := new(MockBcryptUtils)
	userService := services.NewUserService(mockUserRepo, mockBcrypt)

	t.Run("success", func(t *testing.T) {
		req := dto.UserCreateRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"

		mockUserRepo.On("GetUserByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockBcrypt.On("GenerateFromPassword", req.Password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).
			Return(nil).
			Run(func(args mock.Arguments) {
				userArg := args.Get(0).(*domain.User)
				userArg.ID = uuid.New()
			}).Once()

		createdUser, err := userService.CreateUser(req)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.NotEqual(t, uuid.Nil, createdUser.ID)
		require.Equal(t, req.Name, createdUser.Name)
		require.Equal(t, req.Email, createdUser.Email)
		require.Equal(t, hashedPassword, createdUser.Password)
		mockUserRepo.AssertExpectations(t)
		mockBcrypt.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		req := dto.UserCreateRequest{Email: "exists@example.com"}
		existingUser := &domain.User{Email: req.Email}

		mockUserRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil).Once()

		createdUser, err := userService.CreateUser(req)

		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
		require.Nil(t, createdUser)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("bcrypt fails", func(t *testing.T) {
		req := dto.UserCreateRequest{Email: "test@example.com", Password: "password123"}
		bcryptError := errors.New("bcrypt error")

		mockUserRepo.On("GetUserByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockBcrypt.On("GenerateFromPassword", req.Password).Return("", bcryptError).Once()

		createdUser, err := userService.CreateUser(req)

		require.Error(t, err)
		require.Equal(t, bcryptError, err)
		require.Nil(t, createdUser)
		mockUserRepo.AssertExpectations(t)
		mockBcrypt.AssertExpectations(t)
	})

	t.Run("repository create fails", func(t *testing.T) {
		req := dto.UserCreateRequest{Email: "test@example.com", Password: "password12p3"}
		repoError := errors.New("repo error")
		hashedPassword := "hashed_password"

		mockUserRepo.On("GetUserByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockBcrypt.On("GenerateFromPassword", req.Password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(repoError).Once()

		createdUser, err := userService.CreateUser(req)

		require.Error(t, err)
		require.Equal(t, repoError, err)
		require.Nil(t, createdUser)
		mockUserRepo.AssertExpectations(t)
		mockBcrypt.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBcrypt := new(MockBcryptUtils)
	userService := services.NewUserService(mockUserRepo, mockBcrypt)

	t.Run("success", func(t *testing.T) {
		email := "found@example.com"
		expectedUser := &domain.User{ID: uuid.New(), Email: email}
		mockUserRepo.On("GetUserByEmail", email).Return(expectedUser, nil).Once()

		user, err := userService.GetUserByEmail(email)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		email := "notfound@example.com"
		mockUserRepo.On("GetUserByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

		user, err := userService.GetUserByEmail(email)

		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
		require.Nil(t, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("generic db error", func(t *testing.T) {
		email := "db-error@example.com"
		dbErr := errors.New("some database error")
		mockUserRepo.On("GetUserByEmail", email).Return(nil, dbErr).Once()

		user, err := userService.GetUserByEmail(email)

		require.Error(t, err)
		require.Equal(t, dbErr, err)
		require.Nil(t, user)
		mockUserRepo.AssertExpectations(t)
	})
}
