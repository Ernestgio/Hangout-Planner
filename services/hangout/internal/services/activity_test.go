package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type MockActivityRepository struct {
	mock.Mock
}

func (m *MockActivityRepository) WithTx(tx *gorm.DB) repository.ActivityRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.ActivityRepository)
}

func (m *MockActivityRepository) CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	args := m.Called(ctx, activity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Activity), args.Error(1)
}

func (m *MockActivityRepository) GetActivityByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Activity, int64, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).(*domain.Activity), args.Get(1).(int64), args.Error(2)
}

func (m *MockActivityRepository) GetActivitiesByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Activity, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Activity), args.Error(1)
}

func (m *MockActivityRepository) GetAllActivities(ctx context.Context, userID uuid.UUID) ([]repository.ActivityWithCount, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.ActivityWithCount), args.Error(1)
}

func (m *MockActivityRepository) UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	args := m.Called(ctx, activity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Activity), args.Error(1)
}

func (m *MockActivityRepository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestActivityService_CreateActivity(t *testing.T) {
	ctx := context.Background()
	req := &dto.CreateActivityRequest{Name: "Hiking"}
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockActivityRepository)
		checkResult func(t *testing.T, res *dto.ActivityDetailResponse, err error)
	}{
		{
			name: "success",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("CreateActivity", ctx, mock.MatchedBy(func(a *domain.Activity) bool {
					return a.Name == "Hiking" && a.UserID != nil && *a.UserID == userID
				})).
					Return(&domain.Activity{ID: uuid.New(), Name: "Hiking"}, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Hiking", res.Name)
			},
		},
		{
			name: "repository error",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("CreateActivity", ctx, mock.AnythingOfType("*domain.Activity")).
					Return(nil, dbError).Once()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockActivityRepo := new(MockActivityRepository)
			service := services.NewActivityService(nil, mockActivityRepo)
			tc.setupMock(mockActivityRepo)

			result, err := service.CreateActivity(ctx, userID, req)
			tc.checkResult(t, result, err)
			mockActivityRepo.AssertExpectations(t)
		})
	}
}

func TestActivityService_GetActivityByID(t *testing.T) {
	ctx := context.Background()
	activityID := uuid.New()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockActivityRepository)
		checkResult func(t *testing.T, res *dto.ActivityDetailResponse, err error)
	}{
		{
			name: "success",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(&domain.Activity{ID: activityID, Name: "Hiking"}, int64(5), nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, activityID, res.ID)
				require.Equal(t, int64(5), res.HangoutCount)
			},
		},
		{
			name: "not found",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, gorm.ErrRecordNotFound).Once()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrNotFound)
				require.Nil(t, res)
			},
		},
		{
			name: "database error",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, dbError).Once()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockActivityRepository)
			service := services.NewActivityService(nil, mockRepo)
			tc.setupMock(mockRepo)

			result, err := service.GetActivityByID(ctx, activityID, userID)
			tc.checkResult(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestActivityService_GetAllActivities(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockActivityRepository)
		checkResult func(t *testing.T, res []dto.ActivityListItemResponse, err error)
	}{
		{
			name: "success",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetAllActivities", ctx, userID).Return([]repository.ActivityWithCount{
					{Activity: domain.Activity{ID: uuid.New(), Name: "Hiking"}, HangoutCount: 2},
					{Activity: domain.Activity{ID: uuid.New(), Name: "Reading"}, HangoutCount: 5},
				}, nil).Once()
			},
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res, 2)
				require.Equal(t, "Hiking", res[0].Name)
			},
		},
		{
			name: "empty result",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetAllActivities", ctx, userID).Return([]repository.ActivityWithCount{}, nil).Once()
			},
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res, 0)
			},
		},
		{
			name: "database error",
			setupMock: func(repo *MockActivityRepository) {
				repo.On("GetAllActivities", ctx, userID).Return(nil, dbError).Once()
			},
			checkResult: func(t *testing.T, res []dto.ActivityListItemResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockActivityRepository)
			service := services.NewActivityService(nil, mockRepo)
			tc.setupMock(mockRepo)

			result, err := service.GetAllActivities(ctx, userID)
			tc.checkResult(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestActivityService_UpdateActivity(t *testing.T) {
	ctx := context.Background()
	activityID := uuid.New()
	userID := uuid.New()
	req := &dto.UpdateActivityRequest{Name: "Updated Name"}
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock)
		checkResult func(t *testing.T, res *dto.ActivityDetailResponse, err error)
	}{
		{
			name: "success",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(&domain.Activity{ID: activityID, Name: "Old Name"}, int64(2), nil).Once()

				repo.On("UpdateActivity", ctx, mock.MatchedBy(func(act *domain.Activity) bool {
					return act.ID == activityID && act.Name == "Updated Name"
				})).
					Return(&domain.Activity{ID: activityID, Name: "Updated Name"}, nil).Once()

				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Updated Name", res.Name)
				require.Equal(t, int64(2), res.HangoutCount)
			},
		},
		{
			name: "not found",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, gorm.ErrRecordNotFound).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrNotFound)
				require.Nil(t, res)
			},
		},
		{
			name: "db error on get",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, dbError).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "update fails",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(&domain.Activity{ID: activityID, Name: "Old Name"}, int64(3), nil).Once()

				repo.On("UpdateActivity", ctx, mock.AnythingOfType("*domain.Activity")).
					Return(nil, dbError).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.ActivityDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockRepo := new(MockActivityRepository)
			service := services.NewActivityService(db, mockRepo)

			tc.setupMock(mockRepo, sqlMock)
			result, err := service.UpdateActivity(ctx, activityID, userID, req)

			tc.checkResult(t, result, err)
			mockRepo.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestActivityService_DeleteActivity(t *testing.T) {
	ctx := context.Background()
	activityID := uuid.New()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(&domain.Activity{ID: activityID}, int64(0), nil).Once()

				repo.On("DeleteActivity", ctx, activityID).
					Return(nil).Once()

				sqlMock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "not found",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, gorm.ErrRecordNotFound).Once()

				sqlMock.ExpectRollback()
			},
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "get fails",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(nil, 0, dbError).Once()

				sqlMock.ExpectRollback()
			},
			expectedErr: dbError,
		},
		{
			name: "delete fails",
			setupMock: func(repo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo).Once()

				repo.On("GetActivityByID", ctx, activityID, userID).
					Return(&domain.Activity{ID: activityID}, int64(0), nil).Once()

				repo.On("DeleteActivity", ctx, activityID).
					Return(dbError).Once()

				sqlMock.ExpectRollback()
			},
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockRepo := new(MockActivityRepository)
			service := services.NewActivityService(db, mockRepo)

			tc.setupMock(mockRepo, sqlMock)
			err := service.DeleteActivity(ctx, activityID, userID)

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
