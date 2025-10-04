package repository_test

import (
	"database/sql/driver"
	"errors"
	"fmt"
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

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

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

			result, err := repo.CreateHangout(hangout)

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

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		checkResult func(t *testing.T, result *domain.Hangout, err error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				hangoutRows := sqlmock.NewRows([]string{"id", "title", "user_id"}).AddRow(hangoutID, "Found Hangout", &userID)
				userRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(userID, "Test User")

				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(hangoutID, 1).WillReturnRows(hangoutRows)
				mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL").
					WithArgs(userID).WillReturnRows(userRows)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, hangoutID, result.ID)
				require.NotNil(t, result.User)
				require.Equal(t, userID, result.User.ID)
			},
		},
		{
			name: "not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(hangoutID, 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			result, err := repo.GetHangoutByID(hangoutID)
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

			result, err := repo.UpdateHangout(hangout)

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

			err := repo.DeleteHangout(hangoutID)

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

	testCases := []struct {
		name           string
		pagination     *dto.CursorPagination
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult []domain.Hangout
		expectError    bool
	}{
		{
			name:       "first page default sort",
			pagination: &dto.CursorPagination{},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 1")
				expectedSQL := fmt.Sprintf("SELECT * FROM `hangouts` WHERE user_id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY %s %s, id %s LIMIT ?", constants.SortByCreatedAt, constants.SortDirectionDesc, constants.SortDirectionDesc)
				mock.ExpectQuery(expectedSQL).WithArgs(userID, constants.DefaultLimit).WillReturnRows(rows)
			},
			expectedResult: []domain.Hangout{{Title: "Hangout 1"}},
			expectError:    false,
		},
		{
			name:       "second page with cursor sorted by date asc",
			pagination: &dto.CursorPagination{AfterID: &afterID, SortBy: constants.SortByDate, SortDir: string(constants.SortDirectionAsc)},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 2")
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND date >= (SELECT `date` FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL) AND id NOT IN (SELECT `id` FROM `hangouts` WHERE date = (SELECT `date` FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL) AND `hangouts`.`deleted_at` IS NULL) AND `hangouts`.`deleted_at` IS NULL ORDER BY date asc, id asc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, afterID, afterID, constants.DefaultLimit).WillReturnRows(rows)
			},
			expectedResult: []domain.Hangout{{Title: "Hangout 2"}},
			expectError:    false,
		},
		{
			name:       "second page with cursor sorted by created_at desc",
			pagination: &dto.CursorPagination{AfterID: &afterID, SortBy: constants.SortByCreatedAt, SortDir: string(constants.SortDirectionDesc)},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Hangout 3")
				expectedSQL := "SELECT * FROM `hangouts` WHERE user_id = ? AND created_at <= (SELECT `created_at` FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL) AND id NOT IN (SELECT `id` FROM `hangouts` WHERE created_at = (SELECT `created_at` FROM `hangouts` WHERE id = ? AND `hangouts`.`deleted_at` IS NULL) AND `hangouts`.`deleted_at` IS NULL) AND `hangouts`.`deleted_at` IS NULL ORDER BY created_at desc, id desc LIMIT ?"
				mock.ExpectQuery(expectedSQL).WithArgs(userID, afterID, afterID, constants.DefaultLimit).WillReturnRows(rows)
			},
			expectedResult: []domain.Hangout{{Title: "Hangout 3"}},
			expectError:    false,
		},
		{
			name:       "database error",
			pagination: &dto.CursorPagination{},
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedSQL := fmt.Sprintf("SELECT * FROM `hangouts` WHERE user_id = ? AND `hangouts`.`deleted_at` IS NULL ORDER BY %s %s, id %s LIMIT ?", constants.SortByCreatedAt, constants.SortDirectionDesc, constants.SortDirectionDesc)
				mock.ExpectQuery(expectedSQL).WithArgs(userID, constants.DefaultLimit).WillReturnError(errors.New("db error"))
			},
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			results, err := repo.GetHangoutsByUserID(userID, tc.pagination)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, results)
				require.Equal(t, len(tc.expectedResult), len(results))
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
