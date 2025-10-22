package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

type MockHangoutRepository struct {
	mock.Mock
}

func (m *MockHangoutRepository) WithTx(tx *gorm.DB) repository.HangoutRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.HangoutRepository)
}

func (m *MockHangoutRepository) CreateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	args := m.Called(ctx, hangout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hangout), args.Error(1)
}

func (m *MockHangoutRepository) GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Hangout, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hangout), args.Error(1)
}

func (m *MockHangoutRepository) UpdateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	args := m.Called(ctx, hangout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hangout), args.Error(1)
}

func (m *MockHangoutRepository) DeleteHangout(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHangoutRepository) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error) {
	args := m.Called(ctx, userID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Hangout), args.Error(1)
}

func TestHangoutService_CreateHangout(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	validTimeStr := "2025-10-05 15:00:00.000"
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		request     *dto.CreateHangoutRequest
		setupMock   func(repo *MockHangoutRepository)
		checkResult func(t *testing.T, res *dto.HangoutDetailResponse, err error)
	}{
		{
			name: "success",
			request: &dto.CreateHangoutRequest{
				Title: "Test Hangout", Date: validTimeStr, Status: enums.StatusPlanning,
			},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("CreateHangout", ctx, mock.AnythingOfType("*domain.Hangout")).Return(&domain.Hangout{ID: uuid.New(), Title: "Test Hangout"}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Test Hangout", res.Title)
			},
		},
		{
			name:    "success with default status",
			request: &dto.CreateHangoutRequest{Title: "Test Hangout", Date: validTimeStr},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("CreateHangout", ctx, mock.MatchedBy(func(h *domain.Hangout) bool { return h.Status == enums.StatusPlanning })).Return(&domain.Hangout{}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:    "mapper fails on invalid date",
			request: &dto.CreateHangoutRequest{Date: "invalid-date"},
			setupMock: func(repo *MockHangoutRepository) {
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		{
			name:    "repository fails",
			request: &dto.CreateHangoutRequest{Date: validTimeStr},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("CreateHangout", ctx, mock.AnythingOfType("*domain.Hangout")).Return(nil, dbError).Once()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHangoutRepo := new(MockHangoutRepository)
			hangoutService := services.NewHangoutService(nil, mockHangoutRepo)
			tc.setupMock(mockHangoutRepo)

			result, err := hangoutService.CreateHangout(ctx, userID, tc.request)
			tc.checkResult(t, result, err)
			mockHangoutRepo.AssertExpectations(t)
		})
	}
}

func TestHangoutService_GetHangoutByID(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	ownerID := uuid.New()

	testCases := []struct {
		name        string
		userID      uuid.UUID
		setupMock   func(repo *MockHangoutRepository)
		checkResult func(t *testing.T, res *dto.HangoutDetailResponse, err error)
	}{
		{
			name:   "success",
			userID: ownerID,
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("GetHangoutByID", ctx, hangoutID, ownerID).Return(&domain.Hangout{ID: hangoutID, UserID: &ownerID}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, hangoutID, res.ID)
			},
		},
		{
			name:   "not found",
			userID: ownerID,
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("GetHangoutByID", ctx, hangoutID, ownerID).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHangoutRepo := new(MockHangoutRepository)
			hangoutService := services.NewHangoutService(nil, mockHangoutRepo)
			tc.setupMock(mockHangoutRepo)

			result, err := hangoutService.GetHangoutByID(ctx, hangoutID, tc.userID)
			tc.checkResult(t, result, err)
			mockHangoutRepo.AssertExpectations(t)
		})
	}
}

func TestHangoutService_UpdateHangout(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	userID := uuid.New()
	newTitle := "New Title"
	newDateStr := "2025-12-25 18:00:00.000"
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		request     *dto.UpdateHangoutRequest
		setupMock   func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock)
		checkResult func(t *testing.T, res *dto.HangoutDetailResponse, err error)
	}{
		{
			name: "success",
			request: &dto.UpdateHangoutRequest{
				Title:  newTitle,
				Date:   newDateStr,
				Status: enums.StatusConfirmed,
			},
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old Title"}, nil).Once()
				repo.On("UpdateHangout", ctx, mock.MatchedBy(func(h *domain.Hangout) bool { return h.Title == newTitle })).Return(&domain.Hangout{ID: hangoutID, Title: newTitle}, nil).Once()
				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, newTitle, res.Title)
			},
		},
		{
			name: "hangout not found",
			request: &dto.UpdateHangoutRequest{
				Title:  newTitle,
				Date:   newDateStr,
				Status: enums.StatusConfirmed,
			},
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, gorm.ErrRecordNotFound).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrNotFound)
			},
		},
		{
			name: "get hangout by id returns generic db error",
			request: &dto.UpdateHangoutRequest{
				Title:  newTitle,
				Date:   newDateStr,
				Status: enums.StatusConfirmed,
			},
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, dbError).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
		{
			name: "mapper fails on invalid date format",
			request: &dto.UpdateHangoutRequest{
				Title: newTitle,
				Date:  "invalid-date",
			},
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID, UserID: &userID}, nil).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				var parseErr *time.ParseError
				require.ErrorAs(t, err, &parseErr)
			},
		},
		{
			name: "update fails",
			request: &dto.UpdateHangoutRequest{
				Title:  newTitle,
				Date:   newDateStr,
				Status: enums.StatusConfirmed,
			},
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID, UserID: &userID}, nil).Once()
				repo.On("UpdateHangout", ctx, mock.AnythingOfType("*domain.Hangout")).Return(nil, dbError).Once()
				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockRepo := new(MockHangoutRepository)
			service := services.NewHangoutService(db, mockRepo)
			tc.setupMock(mockRepo, sqlMock)

			result, err := service.UpdateHangout(ctx, hangoutID, userID, tc.request)
			tc.checkResult(t, result, err)
			mockRepo.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestHangoutService_DeleteHangout(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID, UserID: &userID}, nil).Once()
				repo.On("DeleteHangout", ctx, hangoutID).Return(nil).Once()
				sqlMock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "hangout not found",
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, gorm.ErrRecordNotFound).Once()
				sqlMock.ExpectRollback()
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "delete fails",
			setupMock: func(repo *MockHangoutRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()
				repo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID, UserID: &userID}, nil).Once()
				repo.On("DeleteHangout", ctx, hangoutID).Return(dbError).Once()
				sqlMock.ExpectRollback()
			},
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockRepo := new(MockHangoutRepository)
			service := services.NewHangoutService(db, mockRepo)
			tc.setupMock(mockRepo, sqlMock)

			err := service.DeleteHangout(ctx, hangoutID, userID)

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestHangoutService_GetHangoutsByUserID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		pagination  *dto.CursorPagination
		setupMock   func(repo *MockHangoutRepository)
		checkResult func(t *testing.T, res *dto.PaginatedHangouts, err error)
	}{
		{
			name:       "success - first page, less results than limit",
			pagination: &dto.CursorPagination{Limit: 10},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return([]domain.Hangout{{ID: uuid.New()}, {ID: uuid.New()}}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Data, 2)
				require.False(t, res.HasMore)
				require.Nil(t, res.NextCursor)
			},
		},
		{
			name:       "success - first page, has more results",
			pagination: &dto.CursorPagination{Limit: 10},
			setupMock: func(repo *MockHangoutRepository) {
				mockHangouts := make([]domain.Hangout, 11) // Fetch limit + 1
				for i := range mockHangouts {
					mockHangouts[i] = domain.Hangout{ID: uuid.New()}
				}
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return(mockHangouts, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Data, 10) // Should be trimmed to limit
				require.True(t, res.HasMore)
				require.NotNil(t, res.NextCursor)
				// Cursor should be the ID of the 10th item
				require.Equal(t, res.Data[9].ID, *res.NextCursor)
			},
		},
		{
			name:       "success - last page",
			pagination: &dto.CursorPagination{Limit: 10},
			setupMock: func(repo *MockHangoutRepository) {
				mockHangouts := make([]domain.Hangout, 10) // Fetch exactly limit
				for i := range mockHangouts {
					mockHangouts[i] = domain.Hangout{ID: uuid.New()}
				}
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return(mockHangouts, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Data, 10)
				require.False(t, res.HasMore)
				require.Nil(t, res.NextCursor)
			},
		},
		{
			name:       "success - no results",
			pagination: &dto.CursorPagination{Limit: 10},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return([]domain.Hangout{}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Data, 0)
				require.False(t, res.HasMore)
				require.Nil(t, res.NextCursor)
			},
		},
		{
			name:       "repository error",
			pagination: &dto.CursorPagination{},
			setupMock: func(repo *MockHangoutRepository) {
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return(nil, dbError).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHangoutRepo := new(MockHangoutRepository)
			hangoutService := services.NewHangoutService(nil, mockHangoutRepo)
			tc.setupMock(mockHangoutRepo)

			result, err := hangoutService.GetHangoutsByUserID(ctx, userID, tc.pagination)
			tc.checkResult(t, result, err)
			mockHangoutRepo.AssertExpectations(t)
		})
	}
}
