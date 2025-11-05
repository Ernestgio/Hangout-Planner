package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewActivityRepository(t *testing.T) {
	db, _ := setupDB(t)
	repo := repository.NewActivityRepository(db)
	require.NotNil(t, repo)
}

func TestActivityRepository_WithTx(t *testing.T) {
	db, _ := setupDB(t)
	repo := repository.NewActivityRepository(db)

	tx := db.Begin()
	txRepo := repo.WithTx(tx)
	require.NotNil(t, txRepo)
	require.NotEqual(t, repo, txRepo)
}

func TestActivityRepository_CreateActivity(t *testing.T) {
	activity := &domain.Activity{Name: "Hiking"}
	dbError := errors.New("create failed")
	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `activities` (`id`,`name`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), activity.Name, AnyTime{}, AnyTime{}, nil, nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `activities` (`id`,`name`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), activity.Name, AnyTime{}, AnyTime{}, nil, nil).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError: true,
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewActivityRepository(db)
			tc.setupMock(mock)

			result, err := repo.CreateActivity(ctx, activity)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEqual(t, uuid.Nil, result.ID)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestActivityRepository_GetActivityByID(t *testing.T) {
	activityID := uuid.New()
	userID := uuid.New()
	ctx := context.Background()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		checkResult func(t *testing.T, act *domain.Activity, count int64, err error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "hangout_count"}).
					AddRow(activityID, "Hiking", 5)
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.id = ? AND activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY `activities`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(activityID, userID, 1).
					WillReturnRows(rows)
			},
			checkResult: func(t *testing.T, act *domain.Activity, count int64, err error) {
				require.NoError(t, err)
				require.NotNil(t, act)
				require.Equal(t, activityID, act.ID)
				require.Equal(t, "Hiking", act.Name)
				require.Equal(t, int64(5), count)
			},
		},
		{
			name: "not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.id = ? AND activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY `activities`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(activityID, userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			checkResult: func(t *testing.T, act *domain.Activity, count int64, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
				require.Nil(t, act)
				require.Zero(t, count)
			},
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.id = ? AND activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY `activities`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(activityID, userID, 1).
					WillReturnError(dbError)
			},
			checkResult: func(t *testing.T, act *domain.Activity, count int64, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, act)
				require.Zero(t, count)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewActivityRepository(db)
			tc.setupMock(mock)

			act, count, err := repo.GetActivityByID(ctx, activityID, userID)
			tc.checkResult(t, act, count, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestActivityRepository_GetAllActivities(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		checkResult func(t *testing.T, results []repository.ActivityWithCount, err error)
	}{
		{
			name: "success with results",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "hangout_count"}).
					AddRow(uuid.New(), "Hiking", 3).
					AddRow(uuid.New(), "Movies", 10)
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY activities.name asc"
				mock.ExpectQuery(expectedSQL).WithArgs(userID).WillReturnRows(rows)
			},
			checkResult: func(t *testing.T, results []repository.ActivityWithCount, err error) {
				require.NoError(t, err)
				require.NotNil(t, results)
				require.Len(t, results, 2)
			},
		},
		{
			name: "success with no results",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "hangout_count"})
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY activities.name asc"
				mock.ExpectQuery(expectedSQL).WithArgs(userID).WillReturnRows(rows)
			},
			checkResult: func(t *testing.T, results []repository.ActivityWithCount, err error) {
				require.NoError(t, err)
				require.NotNil(t, results)
				require.Len(t, results, 0)
			},
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT activities.*, COUNT(hangout_activities.hangout_id) as hangout_count FROM `activities` LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id WHERE activities.user_id = ? AND `activities`.`deleted_at` IS NULL GROUP BY `activities`.`id` ORDER BY activities.name asc"
				mock.ExpectQuery(expectedSQL).WithArgs(userID).WillReturnError(dbError)
			},
			checkResult: func(t *testing.T, results []repository.ActivityWithCount, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, results)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewActivityRepository(db)
			tc.setupMock(mock)

			results, err := repo.GetAllActivities(ctx, userID)
			tc.checkResult(t, results, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestActivityRepository_UpdateActivity(t *testing.T) {
	activity := &domain.Activity{ID: uuid.New(), Name: "Updated Name"}
	activity.CreatedAt = time.Now().Add(-time.Hour)
	dbError := errors.New("update failed")
	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `activities` SET `name`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=? WHERE `activities`.`deleted_at` IS NULL AND `id` = ?").
					WithArgs(activity.Name, activity.CreatedAt, AnyTime{}, nil, nil, activity.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `activities` SET `name`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=? WHERE `activities`.`deleted_at` IS NULL AND `id` = ?").
					WithArgs(activity.Name, activity.CreatedAt, AnyTime{}, nil, nil, activity.ID).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError: true,
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewActivityRepository(db)
			tc.setupMock(mock)

			originalUpdatedAt := activity.UpdatedAt
			result, err := repo.UpdateActivity(ctx, activity)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEqual(t, originalUpdatedAt, result.UpdatedAt)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestActivityRepository_DeleteActivity(t *testing.T) {
	activityID := uuid.New()
	dbError := errors.New("delete failed")
	ctx := context.Background()

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `activities` SET `deleted_at`=? WHERE id = ? AND `activities`.`deleted_at` IS NULL").
					WithArgs(AnyTime{}, activityID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `activities` SET `deleted_at`=? WHERE id = ? AND `activities`.`deleted_at` IS NULL").
					WithArgs(AnyTime{}, activityID).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError: true,
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewActivityRepository(db)
			tc.setupMock(mock)

			err := repo.DeleteActivity(ctx, activityID)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
