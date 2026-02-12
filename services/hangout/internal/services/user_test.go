package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	ctx := context.Background()
	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	dbError := errors.New("db error")
	bcryptError := errors.New("bcrypt error")

	testCases := []struct {
		name        string
		setupMocks  func(repo *MockUserRepository, bcrypt *MockBcryptUtils, sqlMock sqlmock.Sqlmock)
		checkResult func(t *testing.T, user *domain.User, err error)
	}{
		{
			name: "success",
			setupMocks: func(repo *MockUserRepository, bcrypt *MockBcryptUtils, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetUserByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
				bcrypt.On("GenerateFromPassword", req.Password).Return("hashed", nil).Once()
				repo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
					userArg := args.Get(1).(*domain.User)
					userArg.ID = uuid.New()
				}).Return(nil).Once()
				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.NotEqual(t, uuid.Nil, user.ID)
			},
		},
		{
			name: "user already exists",
			setupMocks: func(repo *MockUserRepository, bcrypt *MockBcryptUtils, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetUserByEmail", ctx, req.Email).Return(&domain.User{}, nil).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
			},
		},
		{
			name: "bcrypt fails",
			setupMocks: func(repo *MockUserRepository, bcrypt *MockBcryptUtils, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetUserByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
				bcrypt.On("GenerateFromPassword", req.Password).Return("", bcryptError).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, bcryptError, err)
			},
		},
		{
			name: "repository create fails",
			setupMocks: func(repo *MockUserRepository, bcrypt *MockBcryptUtils, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetUserByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
				bcrypt.On("GenerateFromPassword", req.Password).Return("hashed", nil).Once()
				repo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(dbError).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockRepo := new(MockUserRepository)
			mockBcrypt := new(MockBcryptUtils)
			service := services.NewUserService(db, mockRepo, mockBcrypt, nil)
			tc.setupMocks(mockRepo, mockBcrypt, sqlMock)

			user, err := service.CreateUser(ctx, req)
			tc.checkResult(t, user, err)

			mockRepo.AssertExpectations(t)
			mockBcrypt.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockUserRepository)
		checkResult func(t *testing.T, user *domain.User, err error)
	}{
		{
			name: "success",
			setupMock: func(repo *MockUserRepository) {
				repo.On("GetUserByEmail", ctx, email).Return(&domain.User{Email: email}, nil).Once()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, email, user.Email)
			},
		},
		{
			name: "not found",
			setupMock: func(repo *MockUserRepository) {
				repo.On("GetUserByEmail", ctx, email).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
				require.Nil(t, user)
			},
		},
		{
			name: "database error",
			setupMock: func(repo *MockUserRepository) {
				repo.On("GetUserByEmail", ctx, email).Return(nil, dbError).Once()
			},
			checkResult: func(t *testing.T, user *domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := services.NewUserService(nil, mockRepo, nil, nil)
			tc.setupMock(mockRepo)

			user, err := service.GetUserByEmail(ctx, email)
			tc.checkResult(t, user, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
