package services_test

import (
	"errors"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Simple testify-based mock
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(request dto.UserCreateRequest) (*models.User, error) {
	args := m.Called(request)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAuthService_SignUser(t *testing.T) {
	tests := map[string]struct {
		setupMock func(m *MockUserService)
		input     *dto.SignUpRequest
		wantErr   string
		wantUser  *models.User
	}{
		"user creation fails": {
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything).
					Return(nil, errors.New("db error"))
			},
			input: &dto.SignUpRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Password: "secret",
			},
			wantErr:  "db error",
			wantUser: nil,
		},
		"user creation succeeds": {
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything).
					Return(&models.User{Email: "bob@example.com"}, nil)
			},
			input: &dto.SignUpRequest{
				Name:     "Bob",
				Email:    "bob@example.com",
				Password: "password",
			},
			wantErr:  "",
			wantUser: &models.User{Email: "bob@example.com"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockUserSvc := new(MockUserService)
			tt.setupMock(mockUserSvc)

			authSvc := services.NewAuthService(mockUserSvc)
			user, err := authSvc.SignUser(tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tt.wantUser.Email, user.Email)
			}

			mockUserSvc.AssertExpectations(t)
		})
	}
}
