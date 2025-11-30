package services_test

import (
	"context"
	"errors"
	"testing"

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

func (m *MockHangoutRepository) GetHangoutActivityIDs(ctx context.Context, hangoutID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, hangoutID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockHangoutRepository) AddHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error {
	args := m.Called(ctx, hangoutID, activityIDs)
	return args.Error(0)
}

func (m *MockHangoutRepository) RemoveHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error {
	args := m.Called(ctx, hangoutID, activityIDs)
	return args.Error(0)
}

func TestHangoutService_CreateHangout(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	validTimeStr := "2025-10-05 15:00:00.000"
	dbError := errors.New("db error")
	activityID1 := uuid.New()
	activityID2 := uuid.New()

	type MockSetup func(
		hRepo *MockHangoutRepository,
		aRepo *MockActivityRepository,
		sqlMock sqlmock.Sqlmock,
	)

	testCases := []struct {
		name        string
		request     *dto.CreateHangoutRequest
		setupMock   MockSetup
		checkResult func(t *testing.T, res *dto.HangoutDetailResponse, err error)
	}{
		{
			name: "success_without_activities",
			request: &dto.CreateHangoutRequest{
				Title:  "Test Hangout",
				Date:   validTimeStr,
				Status: enums.StatusPlanning,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				createdHangout := &domain.Hangout{ID: uuid.New(), Title: "Test Hangout"}
				hRepo.On("CreateHangout", ctx, mock.MatchedBy(func(h *domain.Hangout) bool {
					return h.Title == "Test Hangout" && h.Status == enums.StatusPlanning
				})).Return(createdHangout, nil).Once()

				hRepo.On("GetHangoutByID", ctx, createdHangout.ID, userID).Return(createdHangout, nil).Once()

				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Test Hangout", res.Title)
			},
		},
		{
			name: "success_with_empty_activity_slice",
			request: &dto.CreateHangoutRequest{
				Title:       "Empty Activities",
				Date:        validTimeStr,
				ActivityIDs: []uuid.UUID{},
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				createdHangout := &domain.Hangout{ID: uuid.New(), Title: "Empty Activities"}
				hRepo.On("CreateHangout", ctx, mock.MatchedBy(func(h *domain.Hangout) bool {
					return h.Title == "Empty Activities"
				})).Return(createdHangout, nil).Once()

				hRepo.On("GetHangoutByID", ctx, createdHangout.ID, userID).Return(createdHangout, nil).Once()

				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Empty Activities", res.Title)
			},
		},
		{
			name: "success_with_activities",
			request: &dto.CreateHangoutRequest{
				Title:       "Hangout with activities",
				Date:        validTimeStr,
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				mockActivities := []*domain.Activity{{ID: activityID1}, {ID: activityID2}}
				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1, activityID2}).Return(mockActivities, nil).Once()

				createdHangout := &domain.Hangout{ID: uuid.New(), Title: "Hangout with activities"}
				hRepo.On("CreateHangout", ctx, mock.Anything).Return(createdHangout, nil).Once()

				hRepo.On("AddHangoutActivities", ctx, createdHangout.ID, []uuid.UUID{activityID1, activityID2}).Return(nil).Once()

				finalHangout := &domain.Hangout{ID: createdHangout.ID, Title: "Hangout with activities", Activities: []*domain.Activity{{ID: activityID1}, {ID: activityID2}}}
				hRepo.On("GetHangoutByID", ctx, createdHangout.ID, userID).Return(finalHangout, nil).Once()

				sqlMock.ExpectCommit()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Hangout with activities", res.Title)
				require.Len(t, res.Activities, 2)
			},
		},
		{
			name: "mapper_fails_on_invalid_date",
			request: &dto.CreateHangoutRequest{
				Title: "Invalid Date",
				Date:  "invalid-date",
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		{
			name: "activity_validation_fails_invalid_ids",
			request: &dto.CreateHangoutRequest{
				Title:       "Bad IDs",
				Date:        validTimeStr,
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1, activityID2}).
					Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrInvalidActivityIDs)
			},
		},
		{
			name: "activity_repo_error",
			request: &dto.CreateHangoutRequest{
				Title:       "Repo Error",
				Date:        validTimeStr,
				ActivityIDs: []uuid.UUID{activityID1},
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return(nil, dbError).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "repository_create_fails",
			request: &dto.CreateHangoutRequest{
				Title: "DB Error",
				Date:  validTimeStr,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				hRepo.On("CreateHangout", ctx, mock.Anything).Return(nil, dbError).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
		{
			name: "add_activities_fails",
			request: &dto.CreateHangoutRequest{
				Title:       "Add Activity Fail",
				Date:        validTimeStr,
				ActivityIDs: []uuid.UUID{activityID1},
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				mockActivities := []*domain.Activity{{ID: activityID1}}
				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return(mockActivities, nil).Once()

				createdHangout := &domain.Hangout{ID: uuid.New()}
				hRepo.On("CreateHangout", ctx, mock.Anything).Return(createdHangout, nil).Once()

				hRepo.On("AddHangoutActivities", ctx, createdHangout.ID, []uuid.UUID{activityID1}).Return(dbError).Once()

				sqlMock.ExpectRollback()
			},
			checkResult: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
			},
		},
		{
			name: "get_after_create_fails",
			request: &dto.CreateHangoutRequest{
				Title: "After Create Fail",
				Date:  validTimeStr,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				createdHangout := &domain.Hangout{ID: uuid.New()}
				hRepo.On("CreateHangout", ctx, mock.Anything).Return(createdHangout, nil).Once()
				hRepo.On("GetHangoutByID", ctx, createdHangout.ID, userID).Return(nil, dbError).Once()

				sqlMock.ExpectRollback()
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
			db, sqlMock := setupDB(t)
			mockHangoutRepo := new(MockHangoutRepository)
			mockActivityRepo := new(MockActivityRepository)
			service := services.NewHangoutService(db, mockHangoutRepo, mockActivityRepo)

			tc.setupMock(mockHangoutRepo, mockActivityRepo, sqlMock)

			result, err := service.CreateHangout(ctx, userID, tc.request)
			tc.checkResult(t, result, err)

			mockHangoutRepo.AssertExpectations(t)
			mockActivityRepo.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
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
			mockActivityRepo := new(MockActivityRepository)
			hangoutService := services.NewHangoutService(nil, mockHangoutRepo, mockActivityRepo)
			tc.setupMock(mockHangoutRepo)

			result, err := hangoutService.GetHangoutByID(ctx, hangoutID, tc.userID)
			tc.checkResult(t, result, err)
			mockHangoutRepo.AssertExpectations(t)
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
			mockActivityRepo := new(MockActivityRepository)
			service := services.NewHangoutService(db, mockRepo, mockActivityRepo)
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
				mockHangouts := make([]domain.Hangout, 11)
				for i := range mockHangouts {
					mockHangouts[i] = domain.Hangout{ID: uuid.New()}
				}
				repo.On("GetHangoutsByUserID", ctx, userID, mock.AnythingOfType("*dto.CursorPagination")).Return(mockHangouts, nil).Once()
			},
			checkResult: func(t *testing.T, res *dto.PaginatedHangouts, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Data, 10)
				require.True(t, res.HasMore)
				require.NotNil(t, res.NextCursor)
				require.Equal(t, res.Data[9].ID, *res.NextCursor)
			},
		},
		{
			name:       "success - last page",
			pagination: &dto.CursorPagination{Limit: 10},
			setupMock: func(repo *MockHangoutRepository) {
				mockHangouts := make([]domain.Hangout, 10)
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
			mockActivityRepo := new(MockActivityRepository)
			hangoutService := services.NewHangoutService(nil, mockHangoutRepo, mockActivityRepo)
			tc.setupMock(mockHangoutRepo)

			result, err := hangoutService.GetHangoutsByUserID(ctx, userID, tc.pagination)
			tc.checkResult(t, result, err)
			mockHangoutRepo.AssertExpectations(t)
		})
	}
}

func TestHangoutService_UpdateHangout(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	userID := uuid.New()
	dbError := errors.New("db error")
	date := "2025-12-01 18:30:00.000"

	activityID1 := uuid.New()
	activityID2 := uuid.New()

	testCases := []struct {
		name      string
		req       *dto.UpdateHangoutRequest
		setupMock func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock)
		check     func(t *testing.T, res *dto.HangoutDetailResponse, err error)
	}{
		{
			name: "success_no_activity_change",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.MatchedBy(func(h *domain.Hangout) bool {
					return h.ID == hangoutID && h.Title == "Updated Title"
				})).Return(existing, nil).Once()

				final := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Updated Title", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(final, nil).Once()
				sqlMock.ExpectCommit()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Updated Title", res.Title)
				require.Len(t, res.Activities, 1)
			},
		},
		{
			name: "success_add_activity",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1, activityID2}).Return([]*domain.Activity{{ID: activityID1}, {ID: activityID2}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(existing, nil).Once()

				hRepo.On("AddHangoutActivities", ctx, hangoutID, []uuid.UUID{activityID2}).Return(nil).Once()

				final := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Updated Title", Activities: []*domain.Activity{{ID: activityID1}, {ID: activityID2}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(final, nil).Once()
				sqlMock.ExpectCommit()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Activities, 2)
			},
		},
		{
			name: "invalid_activity_ids",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1, activityID2}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, apperrors.ErrInvalidActivityIDs)
				require.Nil(t, res)
			},
		},
		{
			name: "not_found",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, gorm.ErrRecordNotFound).Once()
				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
				require.Nil(t, res)
			},
		},
		{
			name: "remove_activity",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}, {ID: activityID2}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(existing, nil).Once()

				hRepo.On("RemoveHangoutActivities", ctx, hangoutID, []uuid.UUID{activityID2}).Return(nil).Once()

				final := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Updated Title", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(final, nil).Once()
				sqlMock.ExpectCommit()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Activities, 1)
			},
		},
		{
			name: "update_repo_error",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(nil, dbError).Once()
				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "remove_activities_fails",
			req: &dto.UpdateHangoutRequest{
				Title:       "Remove Fail",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}, {ID: activityID2}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(existing, nil).Once()
				hRepo.On("RemoveHangoutActivities", ctx, hangoutID, []uuid.UUID{activityID2}).Return(dbError).Once()
				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "add_activities_fails_update",
			req: &dto.UpdateHangoutRequest{
				Title:       "Add Fail",
				ActivityIDs: []uuid.UUID{activityID1, activityID2},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1, activityID2}).Return([]*domain.Activity{{ID: activityID1}, {ID: activityID2}}, nil).Once()

				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(existing, nil).Once()
				hRepo.On("AddHangoutActivities", ctx, hangoutID, []uuid.UUID{activityID2}).Return(dbError).Once()
				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "final_gethangout_error",
			req: &dto.UpdateHangoutRequest{
				Title:       "Final Get Fail",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()
				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return([]*domain.Activity{{ID: activityID1}}, nil).Once()
				hRepo.On("UpdateHangout", ctx, mock.Anything).Return(existing, nil).Once()
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, dbError).Once()
				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "activity_repo_error",
			req: &dto.UpdateHangoutRequest{
				Title:       "Updated Title",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        date,
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				aRepo.On("GetActivitiesByIDs", ctx, []uuid.UUID{activityID1}).Return(nil, dbError).Once()

				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, res)
			},
		},
		{
			name: "mapper_fails_on_invalid_date",
			req: &dto.UpdateHangoutRequest{
				Title:       "Bad Date",
				ActivityIDs: []uuid.UUID{activityID1},
				Date:        "not-a-date",
			},
			setupMock: func(hRepo *MockHangoutRepository, aRepo *MockActivityRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				hRepo.On("WithTx", mock.Anything).Return(hRepo).Once()
				aRepo.On("WithTx", mock.Anything).Return(aRepo).Once()

				existing := &domain.Hangout{ID: hangoutID, UserID: &userID, Title: "Old", Activities: []*domain.Activity{{ID: activityID1}}}
				hRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(existing, nil).Once()

				sqlMock.ExpectRollback()
			},
			check: func(t *testing.T, res *dto.HangoutDetailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			mockHangoutRepo := new(MockHangoutRepository)
			mockActivityRepo := new(MockActivityRepository)
			service := services.NewHangoutService(db, mockHangoutRepo, mockActivityRepo)

			tc.setupMock(mockHangoutRepo, mockActivityRepo, sqlMock)

			res, err := service.UpdateHangout(ctx, hangoutID, userID, tc.req)
			tc.check(t, res, err)

			mockHangoutRepo.AssertExpectations(t)
			mockActivityRepo.AssertExpectations(t)
			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}
