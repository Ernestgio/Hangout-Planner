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
	"github.com/stretchr/testify/assert"
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
	activityID := uuid.New()
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
				hangoutRows := sqlmock.NewRows([]string{"id", "title", "user_id"}).
					AddRow(hangoutID, "Test Hangout", userID)

				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
					WithArgs(hangoutID, userID, 1).
					WillReturnRows(hangoutRows)

				joinRows := sqlmock.NewRows([]string{"hangout_id", "activity_id"}).
					AddRow(hangoutID, activityID)
				mock.ExpectQuery("SELECT * FROM `hangout_activities` WHERE `hangout_activities`.`hangout_id` = ?").
					WithArgs(hangoutID).
					WillReturnRows(joinRows)

				activityRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(activityID, "Hiking")
				mock.ExpectQuery("SELECT * FROM `activities` WHERE `activities`.`id` = ? AND `activities`.`deleted_at` IS NULL").
					WithArgs(activityID).
					WillReturnRows(activityRows)
			},
			checkResult: func(t *testing.T, result *domain.Hangout, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, hangoutID, result.ID)
				require.Len(t, result.Activities, 1)
				require.Equal(t, activityID, result.Activities[0].ID)
			},
		},
		{
			name: "not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
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
				mock.ExpectQuery("SELECT * FROM `hangouts` WHERE (id = ? AND user_id = ?) AND `hangouts`.`deleted_at` IS NULL ORDER BY `hangouts`.`id` LIMIT ?").
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
	hangoutID := uuid.New()
	ctx := context.Background()
	dbError := errors.New("update error")

	hangoutToUpdate := &domain.Hangout{
		ID:    hangoutID,
		Title: "Updated Title",
		Description: func(s string) *string {
			return &s
		}("New Description"),
	}

	testCases := []struct {
		name        string
		hangout     *domain.Hangout
		setupMock   func(mock sqlmock.Sqlmock, h *domain.Hangout)
		expectError bool
		expectedErr error
	}{
		{
			name:    "success_update",
			hangout: hangoutToUpdate,
			setupMock: func(mock sqlmock.Sqlmock, h *domain.Hangout) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `hangouts` SET `id`=?,`title`=?,`description`=?,`updated_at`=? WHERE id = ? AND `hangouts`.`deleted_at` IS NULL").
					WithArgs(h.ID, h.Title, h.Description, AnyTime{}, h.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:    "database_error",
			hangout: hangoutToUpdate,
			setupMock: func(mock sqlmock.Sqlmock, h *domain.Hangout) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `hangouts` SET `id`=?,`title`=?,`description`=?,`updated_at`=? WHERE id = ? AND `hangouts`.`deleted_at` IS NULL").
					WithArgs(h.ID, h.Title, h.Description, AnyTime{}, h.ID).
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
			tc.setupMock(mock, tc.hangout)

			result, err := repo.UpdateHangout(ctx, tc.hangout)

			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.hangout.ID, result.ID)
				require.Equal(t, tc.hangout.Title, result.Title)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_DeleteHangout(t *testing.T) {
	hangoutID := uuid.New()
	ctx := context.Background()
	dbError := errors.New("db error")

	softDeleteSQL := "UPDATE `hangouts` SET `deleted_at`=? WHERE id = ? AND `hangouts`.`deleted_at` IS NULL"

	testCases := []struct {
		name        string
		id          uuid.UUID
		setupMock   func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectError bool
	}{
		{
			name: "success_full_deletion",
			id:   hangoutID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 5))
				mock.ExpectExec(softDeleteSQL).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "error_on_begin_transaction",
			id:   hangoutID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin().WillReturnError(dbError)
			},
			expectError: true,
		},
		{
			name: "error_on_join_table_delete_triggers_rollback",
			id:   hangoutID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?").
					WithArgs(id).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "error_on_hangout_soft_delete_triggers_rollback",
			id:   hangoutID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 5))
				mock.ExpectExec(softDeleteSQL).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "error_on_commit",
			id:   hangoutID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 5))
				mock.ExpectExec(softDeleteSQL).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(dbError)
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)

			tc.setupMock(mock, tc.id)

			err := repo.DeleteHangout(ctx, tc.id)

			if tc.expectError {
				require.Error(t, err)
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

func TestHangoutRepository_GetHangoutActivityIDs(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	activityID1 := uuid.New()
	activityID2 := uuid.New()
	dbError := errors.New("db error")

	testCases := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedIDs   []uuid.UUID
		expectError   bool
		expectedError error
	}{
		{
			name: "success_fetch",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"activity_id"}).
					AddRow(activityID1.String()).
					AddRow(activityID2.String())
				mock.ExpectQuery("SELECT `activity_id` FROM `hangout_activities` WHERE hangout_id = ?").
					WithArgs(hangoutID).
					WillReturnRows(rows)
			},
			expectedIDs:   []uuid.UUID{activityID1, activityID2},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "no_results",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"activity_id"})
				mock.ExpectQuery("SELECT `activity_id` FROM `hangout_activities` WHERE hangout_id = ?").
					WithArgs(hangoutID).
					WillReturnRows(rows)
			},
			expectedIDs:   []uuid.UUID{},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "database_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT `activity_id` FROM `hangout_activities` WHERE hangout_id = ?").
					WithArgs(hangoutID).
					WillReturnError(dbError)
			},
			expectedIDs:   nil,
			expectError:   true,
			expectedError: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			ids, err := repo.GetHangoutActivityIDs(ctx, hangoutID)

			if tc.expectError {
				require.Error(t, err)
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				}
				require.Nil(t, ids)
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, tc.expectedIDs, ids)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_AddHangoutActivities(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	activityID1 := uuid.New()
	activityID2 := uuid.New()
	activityIDs := []uuid.UUID{activityID1, activityID2}
	dbError := errors.New("insert failed")

	testCases := []struct {
		name          string
		activityIDs   []uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		expectedError error
	}{
		{
			name:        "success_insert_multiple",
			activityIDs: activityIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `hangout_activities` (`activity_id`,`hangout_id`) VALUES (?,?),(?,?)").
					WithArgs(activityID1, hangoutID, activityID2, hangoutID).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:        "empty_list",
			activityIDs: []uuid.UUID{},
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: false,
		},
		{
			name:        "database_error",
			activityIDs: activityIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `hangout_activities` (`activity_id`,`hangout_id`) VALUES (?,?),(?,?)").
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError:   true,
			expectedError: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			err := repo.AddHangoutActivities(ctx, hangoutID, tc.activityIDs)

			if tc.expectError {
				require.Error(t, err)
				if !assert.Equal(t, tc.expectedError, err) {
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHangoutRepository_RemoveHangoutActivities(t *testing.T) {
	ctx := context.Background()
	hangoutID := uuid.New()
	activityID1 := uuid.New()
	activityID2 := uuid.New()
	activityIDs := []uuid.UUID{activityID1, activityID2}
	dbError := errors.New("delete failed")

	testCases := []struct {
		name          string
		activityIDs   []uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		expectedError error
	}{
		{
			name:        "success_delete_multiple",
			activityIDs: activityIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE hangout_id = ? AND activity_id IN (?,?)").
					WithArgs(hangoutID, activityID1, activityID2).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:        "empty_list",
			activityIDs: []uuid.UUID{},
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: false,
		},
		{
			name:        "database_error",
			activityIDs: activityIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `hangout_activities` WHERE hangout_id = ? AND activity_id IN (?,?)").
					WithArgs(hangoutID, activityID1, activityID2).
					WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectError:   true,
			expectedError: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupDB(t)
			repo := repository.NewHangoutRepository(db)
			tc.setupMock(mock)

			err := repo.RemoveHangoutActivities(ctx, hangoutID, tc.activityIDs)

			if tc.expectError {
				require.Error(t, err)
				if !assert.Equal(t, tc.expectedError, err) {
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
