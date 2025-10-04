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
)

func TestAuthService_SignUser(t *testing.T) {
	mockJwtSvc := new(MockJWTUtils)
	mockBcrypt := new(MockBcryptUtils)
	newUserID := uuid.New()

	tests := map[string]struct {
		setupMock func(m *MockUserService)
		input     *dto.SignUpRequest
		wantErr   string
		wantUser  *domain.User
	}{
		"User creation fails": {
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything).
					Return(nil, errors.New("db error"))
			},
			input:    &dto.SignUpRequest{Name: "Alice", Email: "alice@example.com", Password: "secret"},
			wantErr:  "db error",
			wantUser: nil,
		},
		"User creation succeeds": {
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything).
					Return(&domain.User{ID: newUserID, Email: "bob@example.com"}, nil)
			},
			input:    &dto.SignUpRequest{Name: "Bob", Email: "bob@example.com", Password: "password"},
			wantErr:  "",
			wantUser: &domain.User{ID: newUserID, Email: "bob@example.com"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockUserSvc := new(MockUserService)
			tt.setupMock(mockUserSvc)

			authSvc := services.NewAuthService(mockUserSvc, mockJwtSvc, mockBcrypt)
			user, err := authSvc.SignUser(tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tt.wantUser.Email, user.Email)
				require.Equal(t, tt.wantUser.ID, user.ID)
			}

			mockUserSvc.AssertExpectations(t)
		})
	}
}

func TestAuthService_SignInUser(t *testing.T) {
	const correctPassword = "StrongPassword123"
	const correctEmail = "user@valid.com"
	const mockToken = "signed.jwt.token"

	validUserID := uuid.New()
	validUser := &domain.User{
		ID:       validUserID,
		Email:    correctEmail,
		Password: "hashed-password",
	}

	tests := map[string]struct {
		setupUserMock   func(m *MockUserService)
		setupBcryptMock func(m *MockBcryptUtils)
		setupJWTMock    func(m *MockJWTUtils)
		input           *dto.SignInRequest
		wantErr         error
		wantToken       string
	}{
		"Failure_UserNotFound": {
			setupUserMock: func(m *MockUserService) {
				m.On("GetUserByEmail", "notfound@email.com").Return(nil, nil)
			},
			setupBcryptMock: func(m *MockBcryptUtils) {},
			setupJWTMock:    func(m *MockJWTUtils) {},
			input:           &dto.SignInRequest{Email: "notfound@email.com", Password: "any"},
			wantErr:         apperrors.ErrInvalidCredentials,
			wantToken:       "",
		},
		"Failure_UserServiceError": {
			setupUserMock: func(m *MockUserService) {
				m.On("GetUserByEmail", correctEmail).Return(nil, errors.New("db connection error"))
			},
			setupBcryptMock: func(m *MockBcryptUtils) {},
			setupJWTMock:    func(m *MockJWTUtils) {},
			input:           &dto.SignInRequest{Email: correctEmail, Password: "any"},
			wantErr:         errors.New("db connection error"),
			wantToken:       "",
		},
		"Failure_PasswordMismatch": {
			setupUserMock: func(m *MockUserService) {
				m.On("GetUserByEmail", correctEmail).Return(validUser, nil)
			},
			setupBcryptMock: func(m *MockBcryptUtils) {
				m.On("CompareHashAndPassword", validUser.Password, "WrongPassword").
					Return(apperrors.ErrInvalidCredentials)
			},
			setupJWTMock: func(m *MockJWTUtils) {},
			input:        &dto.SignInRequest{Email: correctEmail, Password: "WrongPassword"},
			wantErr:      apperrors.ErrInvalidCredentials,
			wantToken:    "",
		},
		"Failure_JWTGenerationError": {
			setupUserMock: func(m *MockUserService) {
				m.On("GetUserByEmail", correctEmail).Return(validUser, nil)
			},
			setupBcryptMock: func(m *MockBcryptUtils) {
				m.On("CompareHashAndPassword", validUser.Password, correctPassword).Return(nil)
			},
			setupJWTMock: func(m *MockJWTUtils) {
				m.On("Generate", validUser).Return("", errors.New("jwt signing failed"))
			},
			input:     &dto.SignInRequest{Email: correctEmail, Password: correctPassword},
			wantErr:   errors.New("jwt signing failed"),
			wantToken: "",
		},
		"Success_ValidCredentials": {
			setupUserMock: func(m *MockUserService) {
				m.On("GetUserByEmail", correctEmail).Return(validUser, nil)
			},
			setupBcryptMock: func(m *MockBcryptUtils) {
				m.On("CompareHashAndPassword", validUser.Password, correctPassword).Return(nil)
			},
			setupJWTMock: func(m *MockJWTUtils) {
				m.On("Generate", validUser).Return(mockToken, nil)
			},
			input:     &dto.SignInRequest{Email: correctEmail, Password: correctPassword},
			wantErr:   nil,
			wantToken: mockToken,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockUserSvc := new(MockUserService)
			mockJwtSvc := new(MockJWTUtils)
			mockBcrypt := new(MockBcryptUtils)

			tt.setupUserMock(mockUserSvc)
			tt.setupBcryptMock(mockBcrypt)
			tt.setupJWTMock(mockJwtSvc)

			authSvc := services.NewAuthService(mockUserSvc, mockJwtSvc, mockBcrypt)
			response, err := authSvc.SignInUser(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, apperrors.ErrInvalidCredentials) {
					require.ErrorIs(t, err, tt.wantErr)
				} else {
					require.EqualError(t, err, tt.wantErr.Error())
				}
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tt.wantToken, response.Token)
			}

			mockUserSvc.AssertExpectations(t)
			mockJwtSvc.AssertExpectations(t)
			mockBcrypt.AssertExpectations(t)
		})
	}
}
