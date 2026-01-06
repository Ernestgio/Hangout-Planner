package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	repo "github.com/Ernestgio/Hangout-Planner/services/file/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreate_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, *domain.MemoryFile)
		wantError bool
	}{
		{
			name: "success",
			prepare: func(m sqlmock.Sqlmock, f *domain.MemoryFile) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memory_files.*").WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "insert error",
			prepare: func(m sqlmock.Sqlmock, f *domain.MemoryFile) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memory_files.*").WillReturnError(errors.New("insert failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			f := &domain.MemoryFile{OriginalName: "test.png", MemoryID: uuid.New()}
			tt.prepare(mock, f)
			got, err := r.Create(ctx, f)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, f.OriginalName, got.OriginalName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateBatch_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		files     []*domain.MemoryFile
		prepare   func(sqlmock.Sqlmock, []*domain.MemoryFile)
		wantError bool
	}{
		{
			name: "success multiple files",
			files: []*domain.MemoryFile{
				{OriginalName: "a.png", MemoryID: uuid.New()},
				{OriginalName: "b.png", MemoryID: uuid.New()},
			},
			prepare: func(m sqlmock.Sqlmock, f []*domain.MemoryFile) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memory_files.*").WillReturnResult(sqlmock.NewResult(1, 2))
				m.ExpectCommit()
			},
		},
		{
			name: "insert error",
			files: []*domain.MemoryFile{
				{OriginalName: "a.png", MemoryID: uuid.New()},
			},
			prepare: func(m sqlmock.Sqlmock, f []*domain.MemoryFile) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memory_files.*").WillReturnError(errors.New("batch insert failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			tt.prepare(mock, tt.files)
			err := r.CreateBatch(ctx, tt.files)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByID_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, uuid.UUID)
		wantError bool
	}{
		{
			name: "found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				cols := []string{"id", "original_name", "file_extension", "storage_path", "file_size", "mime_type", "file_status", "created_at", "deleted_at", "memory_id"}
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").
					WithArgs(id, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(cols).AddRow(id, "file.png", "png", "/path/file.png", 1024, "image/png", "PENDING", time.Now(), nil, uuid.New()))
			},
		},
		{
			name: "not found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").
					WithArgs(id, sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			id := uuid.New()
			tt.prepare(mock, id)
			f, err := r.GetByID(ctx, id)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "file.png", f.OriginalName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByMemoryID_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, uuid.UUID)
		wantError bool
	}{
		{
			name: "found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				cols := []string{"id", "original_name", "file_extension", "storage_path", "file_size", "mime_type", "file_status", "created_at", "deleted_at", "memory_id"}
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").
					WithArgs(id, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(cols).AddRow(uuid.New(), "mem.png", "png", "/p/mem.png", 2048, "image/png", "UPLOADED", time.Now(), nil, id))
			},
		},
		{
			name: "not found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").
					WithArgs(id, sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			id := uuid.New()
			tt.prepare(mock, id)
			f, err := r.GetByMemoryID(ctx, id)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "mem.png", f.OriginalName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByMemoryIDs_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		ids       []uuid.UUID
		prepare   func(sqlmock.Sqlmock, []uuid.UUID)
		wantCount int
		wantError bool
	}{
		{
			name: "found multiple",
			ids:  []uuid.UUID{uuid.New(), uuid.New()},
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID) {
				cols := []string{"id", "original_name", "file_extension", "storage_path", "file_size", "mime_type", "file_status", "created_at", "deleted_at", "memory_id"}
				rows := sqlmock.NewRows(cols).
					AddRow(uuid.New(), "a.png", "png", "/p/a.png", 100, "image/png", "PENDING", time.Now(), nil, ids[0]).
					AddRow(uuid.New(), "b.png", "png", "/p/b.png", 200, "image/png", "UPLOADED", time.Now(), nil, ids[1])
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "empty result",
			ids:  []uuid.UUID{uuid.New()},
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID) {
				cols := []string{"id", "original_name", "file_extension", "storage_path", "file_size", "mime_type", "file_status", "created_at", "deleted_at", "memory_id"}
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").WillReturnRows(sqlmock.NewRows(cols))
			},
			wantCount: 0,
		},
		{
			name: "query error",
			ids:  []uuid.UUID{uuid.New()},
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").WillReturnError(errors.New("query failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			tt.prepare(mock, tt.ids)
			files, err := r.GetByMemoryIDs(ctx, tt.ids)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, files, tt.wantCount)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateStatusBatch_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		ids       []uuid.UUID
		status    string
		prepare   func(sqlmock.Sqlmock, []uuid.UUID, string)
		wantError bool
	}{
		{
			name:   "success multiple ids",
			ids:    []uuid.UUID{uuid.New(), uuid.New()},
			status: "UPLOADED",
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID, status string) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memory_files.* SET .*file_status.*").
					WithArgs(status, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 2))
				m.ExpectCommit()
			},
		},
		{
			name:   "empty ids",
			ids:    []uuid.UUID{},
			status: "UPLOADED",
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID, status string) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memory_files.* SET .*file_status.*").
					WithArgs(status).
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectCommit()
			},
		},
		{
			name:   "update error",
			ids:    []uuid.UUID{uuid.New()},
			status: "UPLOADED",
			prepare: func(m sqlmock.Sqlmock, ids []uuid.UUID, status string) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memory_files.* SET .*file_status.*").
					WithArgs(status, sqlmock.AnyArg()).
					WillReturnError(errors.New("update failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			tt.prepare(mock, tt.ids, tt.status)
			err := r.UpdateStatusBatch(ctx, tt.ids, tt.status)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDelete_TableDriven(t *testing.T) {
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
				m.ExpectExec("UPDATE .*memory_files.* SET .*deleted_at.*").
					WithArgs(AnyTime{}, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "delete error",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memory_files.* SET .*deleted_at.*").
					WithArgs(AnyTime{}, id).
					WillReturnError(errors.New("delete failed"))
				m.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newDBWithRegexp(t)
			r := repo.NewMemoryFileRepository(db)
			id := uuid.New()
			tt.prepare(mock, id)
			err := r.Delete(ctx, id)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWithTx_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		prepare func(sqlmock.Sqlmock)
	}{
		{
			name: "use tx db",
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO .*memory_files.*").WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMain, mockMain := newDBWithRegexp(t)
			dbTx, mockTx := newDBWithRegexp(t)
			mainRepo := repo.NewMemoryFileRepository(dbMain)
			f := &domain.MemoryFile{OriginalName: "tx.png", MemoryID: uuid.New()}
			tt.prepare(mockTx)
			r := mainRepo.WithTx(dbTx)
			got, err := r.Create(ctx, f)
			require.NoError(t, err)
			require.Equal(t, "tx.png", got.OriginalName)
			require.NoError(t, mockTx.ExpectationsWereMet())
			require.NoError(t, mockMain.ExpectationsWereMet())
		})
	}
}
