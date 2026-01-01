package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	repo "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateMemory_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock)
		wantError bool
	}{
		{
			name: "success",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memories.*").WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "insert error",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memories.*").WillReturnError(errors.New("insert failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryRepository(db)
			m := &domain.Memory{Name: "m", HangoutID: uuid.New(), UserID: uuid.New()}
			tt.prepare(mock)
			got, err := r.CreateMemory(ctx, m)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, m.Name, got.Name)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMemoryByID_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, uuid.UUID, uuid.UUID)
		wantError bool
	}{
		{
			name: "found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID, userID uuid.UUID) {
				cols := []string{"id", "name", "created_at", "updated_at", "deleted_at", "hangout_id", "user_id"}
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(id, userID, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows(cols).AddRow(uuid.New(), "nm", time.Now(), time.Now(), nil, uuid.New(), userID))
			},
		},
		{
			name: "not found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID, userID uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(id, userID, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},
		{
			name: "db error",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID, userID uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(id, userID, sqlmock.AnyArg()).WillReturnError(errors.New("db error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryRepository(db)
			id := uuid.New()
			userID := uuid.New()
			tt.prepare(mock, id, userID)
			mem, err := r.GetMemoryByID(ctx, id, userID)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "nm", mem.Name)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMemoriesByHangoutID_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		pagination *dto.CursorPagination
		prepare    func(sqlmock.Sqlmock, uuid.UUID, *dto.CursorPagination)
		wantError  bool
	}{
		{
			name:       "no cursor",
			pagination: &dto.CursorPagination{Limit: 2, SortDir: "asc"},
			prepare: func(m sqlmock.Sqlmock, hangoutID uuid.UUID, p *dto.CursorPagination) {
				cols := []string{"id", "name", "created_at", "updated_at", "deleted_at", "hangout_id", "user_id"}
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(hangoutID, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows(cols).AddRow(uuid.New(), "a", time.Now(), time.Now(), nil, hangoutID, uuid.New()).AddRow(uuid.New(), "b", time.Now(), time.Now(), nil, hangoutID, uuid.New()))
			},
		},
		{
			name:       "invalid cursor",
			pagination: &dto.CursorPagination{Limit: 2, AfterID: func() *uuid.UUID { id := uuid.New(); return &id }()},
			prepare: func(m sqlmock.Sqlmock, hangoutID uuid.UUID, p *dto.CursorPagination) {
				aid := *p.AfterID
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(aid, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},
		{
			name:       "db error",
			pagination: &dto.CursorPagination{Limit: 1, SortDir: "desc"},
			prepare: func(m sqlmock.Sqlmock, hangoutID uuid.UUID, p *dto.CursorPagination) {
				m.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(hangoutID, sqlmock.AnyArg()).WillReturnError(errors.New("db error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryRepository(db)
			hid := uuid.New()
			tt.prepare(mock, hid, tt.pagination)
			res, err := r.GetMemoriesByHangoutID(ctx, hid, tt.pagination)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(res), 1)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteMemory_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, uuid.UUID)
		wantError bool
	}{
		{
			name: "success",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memories.*SET .*deleted_at.*").WithArgs(AnyTime{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "delete error",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memories.*SET .*deleted_at.*").WithArgs(AnyTime{}, id).WillReturnError(errors.New("delete failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryRepository(db)
			id := uuid.New()
			tt.prepare(mock, id)
			err := r.DeleteMemory(ctx, id)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMemoriesByHangoutID_WithCursor_TableDriven(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		sortDir string
	}{
		{name: "cursor asc", sortDir: "asc"},
		{name: "cursor desc", sortDir: "desc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryRepository(db)

			hid := uuid.New()
			cursorID := uuid.New()
			cursorCreated := time.Now().Add(-time.Hour)

			cols := []string{"id", "name", "created_at", "updated_at", "deleted_at", "hangout_id", "user_id"}

			mock.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(cursorID, sqlmock.AnyArg()).WillReturnRows(
				sqlmock.NewRows(cols).AddRow(cursorID, "c", cursorCreated, cursorCreated, nil, hid, uuid.New()),
			)

			comp1 := cursorCreated
			comp2 := cursorCreated
			mock.ExpectQuery("SELECT .* FROM .*memories.*").WithArgs(hid, comp1, comp2, cursorID, sqlmock.AnyArg()).WillReturnRows(
				sqlmock.NewRows(cols).AddRow(uuid.New(), "r", time.Now(), time.Now(), nil, hid, uuid.New()),
			)

			p := &dto.CursorPagination{Limit: 1, SortDir: tc.sortDir, AfterID: &cursorID}
			res, err := r.GetMemoriesByHangoutID(ctx, hid, p)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(res), 1)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMemoryRepository_WithTx(t *testing.T) {
	ctx := context.Background()
	dbMain, mockMain := newDBWithRegexp(t)
	dbTx, mockTx := newDBWithRegexp(t)

	mainRepo := repo.NewMemoryRepository(dbMain)
	m := &domain.Memory{Name: "txm", HangoutID: uuid.New(), UserID: uuid.New()}

	mockTx.ExpectBegin()
	mockTx.ExpectExec("INSERT INTO .*memories.*").WillReturnResult(sqlmock.NewResult(1, 1))
	mockTx.ExpectCommit()

	r := mainRepo.WithTx(dbTx)
	got, err := r.CreateMemory(ctx, m)
	require.NoError(t, err)
	require.Equal(t, "txm", got.Name)

	require.NoError(t, mockTx.ExpectationsWereMet())
	require.NoError(t, mockMain.ExpectationsWereMet())
}
