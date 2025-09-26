package services_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := map[string]struct {
		setupRepo func(m *MockUserRepository)
		request   dto.UserCreateRequest
		wantErr   string
	}{
		"user already exists": {
			setupRepo: func(m *MockUserRepository) {
				m.On("GetUserByEmail", "exists@example.com").
					Return(&models.User{}, nil)
			},
			request: dto.UserCreateRequest{
				Name:     "Alice",
				Email:    "exists@example.com",
				Password: "password",
			},
			wantErr: apperrors.ErrUserAlreadyExists.Error(),
		},
		"password hashing fails": {
			setupRepo: func(m *MockUserRepository) {
				m.On("GetUserByEmail", "hashfail@example.com").
					Return(nil, errors.New("not found"))
			},
			request: dto.UserCreateRequest{
				Name:     "Bob",
				Email:    "hashfail@example.com",
				Password: string(make([]byte, 1000000000)),
			},
			wantErr: "bcrypt",
		},
		"repo create fails": {
			setupRepo: func(m *MockUserRepository) {
				m.On("GetUserByEmail", "fail@example.com").
					Return(nil, errors.New("not found"))
				m.On("CreateUser", mock.Anything).
					Return(errors.New("db error"))
			},
			request: dto.UserCreateRequest{
				Name:     "Carol",
				Email:    "fail@example.com",
				Password: "password",
			},
			wantErr: "db error",
		},
		"success": {
			setupRepo: func(m *MockUserRepository) {
				m.On("GetUserByEmail", "ok@example.com").
					Return(nil, errors.New("not found"))
				m.On("CreateUser", mock.Anything).
					Return(nil)
			},
			request: dto.UserCreateRequest{
				Name:     "Dave",
				Email:    "ok@example.com",
				Password: "password",
			},
			wantErr: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupRepo(mockRepo)

			svc := services.NewUserService(mockRepo, utils.NewBcryptUtils(4))

			user, err := svc.CreateUser(tt.request)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tt.request.Email, user.Email)
				require.NotEmpty(t, user.Password)
				require.NotEmpty(t, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	tests := map[string]struct {
		setupRepo func(m *MockUserRepository)
		email     string
		wantUser  *models.User
		wantErr   string
	}{
		"success": {
			setupRepo: func(m *MockUserRepository) {
				expected := &models.User{Email: "test@example.com"}
				m.On("GetUserByEmail", "test@example.com").
					Return(expected, nil)
			},
			email:    "test@example.com",
			wantUser: &models.User{Email: "test@example.com"},
			wantErr:  "",
		},
		"repo error": {
			setupRepo: func(m *MockUserRepository) {
				m.On("GetUserByEmail", "fail@example.com").
					Return(nil, errors.New("db error"))
			},
			email:    "fail@example.com",
			wantUser: nil,
			wantErr:  "db error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupRepo(mockRepo)

			svc := services.NewUserService(mockRepo, utils.NewBcryptUtils(4))

			user, err := svc.GetUserByEmail(tt.email)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tt.wantUser.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
