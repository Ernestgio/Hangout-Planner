package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewHangoutRepository(t *testing.T) {
	db, _ := setupDB(t)
	repo := repository.NewHangoutRepository(db)
	require.NotNil(t, repo)
}

func TestHangoutRepository_WithTx(t *testing.T) {
	db, _ := setupDB(t)
	repo := repository.NewHangoutRepository(db)

	tx := db.Begin()
	txRepo := repo.WithTx(tx)
	require.NotNil(t, txRepo)
	require.NotEqual(t, repo, txRepo)
}

func TestHangoutRepository_CreateHangout(t *testing.T) {
	hangout := &domain.Hangout{Title: "Test Hangout"}
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
				mock.ExpectExec("INSERT INTO `hangouts` (`id`,`title`,`description`,`date`,`status`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), hangout.Title, hangout.Description, hangout.Date, hangout.Status, AnyTime{}, AnyTime{}, nil, nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `hangouts` (`id`,`title`,`description`,`date`,`status`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), hangout.Title, hangout.Description, hangout.Date, hangout.Status, AnyTime{}, AnyTime{}, nil, nil).
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
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			result, err := repo.CreateHangout(ctx, hangout)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_GetHangoutByID(t *testing.T) {
	hangoutID := uuid.New()
	userID := uuid.New()
	ctx := context.Background()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		checkResult func(t *testing.T, result *domain.Hangout, err error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "user_id"}).
					AddRow(hangoutID, "Found Hangout", &userID)
				expectedSQL := "SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(hangoutID, userID, 1).
					WillReturnRows(rows)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, hangoutID, result.ID)
				require.Equal(t, &userID, result.UserID)
			},
		},
		{
			name: "not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(hangoutID, userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
				require.Nil(t, result)
			},
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?"
				mock.ExpectQuery(expectedSQL).
					WithArgs(hangoutID, userID, 1).
					WillReturnError(dbError)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.Error(t, err)
				require.Equal(t, dbError, err)
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			result, err := repo.GetHangoutByID(ctx, hangoutID, userID)
			tc.checkResult(t, result, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_UpdateHangout(t *testing.T) {
	hangout := &domain.Hangout{
		ID:    uuid.New(),
		Title: "Updated Title",
	}
	hangout.CreatedAt = time.Now()
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
				mock.ExpectExec("UPDATE `hangouts` SET `title`=?,`description`=?,`date`=?,`status`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=? WHERE `hangouts`.`deleted_at` IS NULL AND `id` = ?").
					WithArgs(hangout.Title, hangout.Description, hangout.Date, hangout.Status, hangout.CreatedAt, AnyTime{}, nil, nil, hangout.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `hangouts` SET `title`=?,`description`=?,`date`=?,`status`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=? WHERE `hangouts`.`deleted_at` IS NULL AND `id` = ?").
					WithArgs(hangout.Title, hangout.Description, hangout.Date, hangout.Status, hangout.CreatedAt, AnyTime{}, nil, nil, hangout.ID).
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
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			result, err := repo.UpdateHangout(ctx, hangout)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, hangout.Title, result.Title)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_DeleteHangout(t *testing.T) {
	hangoutID := uuid.New()
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
				mock.ExpectExec("UPDATE `hangouts` SET `deleted_at`=? WHERE id = ? AND `hangouts`.`deleted_at` IS NULL").
					WithArgs(AnyTime{}, hangoutID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `hangouts` SET `deleted_at`=? WHERE id = ? AND `hangouts`.`deleted_at` IS NULL").
					WithArgs(AnyTime{}, hangoutID).
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
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			err := repo.DeleteHangout(ctx, hangoutID)

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

func TestHangoutRepository_GetHangoutsByUserID(t *testing.T) {
	userID := uuid.New()
	afterID := uuid.New()
	cursorTime := time.Now().Add(-1 * time.Hour)
	ctx := context.Background()
	dbError := errors.New("db error")

	testCases := []struct {
		name        string
		pagination  *dto.CursorPagination
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		expectedLen int
	}{
		{
			name:       "first page default sort (created_at desc)",
			pagination: &dto.CursorPagination{},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 1")
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY created_at desc, id desc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, constants.DefaultLimit+1).WillReturnRows(rows)
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name:       "second page with cursor sorted by date asc",
			pagination: &dto.CursorPagination{AfterID: &afterID, SortBy: constants.SortByDate, SortDir: string(constants.SortDirectionAsc), Limit: 15},
			setupMock: func(mock sqlmock.Sqlmock) {
				cursorRows := sqlmock.NewRows([]string{"id", "date", "created_at"}).AddRow(afterID, cursorTime, cursorTime)
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(afterID, 1).WillReturnRows(cursorRows)

				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 2")
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND ((date > ?) OR (date = ? AND id > ?)) AND `hangouts`.`deleted_at` IS NULL ORDER BY date asc, id asc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, cursorTime, cursorTime, afterID, 15+1).WillReturnRows(rows)
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name:       "second page with cursor sorted by created_at desc",
			pagination: &dto.CursorPagination{AfterID: &afterID, SortBy: constants.SortByCreatedAt, SortDir: string(constants.SortDirectionDesc), Limit: 5},
			setupMock: func(mock sqlmock.Sqlmock) {
				cursorRows := sqlmock.NewRows([]string{"id", "date", "created_at"}).AddRow(afterID, cursorTime, cursorTime)
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(afterID, 1).WillReturnRows(cursorRows)

				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 3")
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND ((created_at < ?) OR (created_at = ? AND id < ?)) AND `hangouts`.`deleted_at` IS NULL ORDER BY created_at desc, id desc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, cursorTime, cursorTime, afterID, 5+1).WillReturnRows(rows)
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name:       "database error on main query",
			pagination: &dto.CursorPagination{},
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY created_at desc, id desc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, constants.DefaultLimit+1).WillReturnError(dbError)
			},
			expectError: true,
			expectedLen: 0,
		},
		{
			name:       "database error on cursor fetch",
			pagination: &dto.CursorPagination{AfterID: &afterID},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(afterID, 1).WillReturnError(dbError)
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			results, err := repo.GetHangoutsByUserID(ctx, userID, tc.pagination)

			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, results)
			} else {
				require.NoError(t, err)
				require.NotNil(t, results)
				require.Len(t, results, tc.expectedLen)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
