package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	repo "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateFile_TableDriven(t *testing.T) {
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
			f := &domain.MemoryFile{OriginalName: "x.png", MemoryID: uuid.New()}
			tt.prepare(mock, f)
			got, err := r.CreateFile(ctx, f)
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

func TestGetFileByMemoryID_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prepare   func(sqlmock.Sqlmock, uuid.UUID)
		wantError bool
	}{
		{
			name: "found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				cols := []string{"id", "original_name", "file_extension", "storage_path", "file_size", "mime_type", "created_at", "deleted_at", "memory_id"}
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").WithArgs(id, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows(cols).AddRow(uuid.New(), "n.png", "png", "/p/n.png", 10, "image/png", time.Now(), nil, id))
			},
		},
		{
			name: "not found",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectQuery("SELECT .* FROM .*memory_files.*").WithArgs(id, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
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
			f, err := r.GetFileByMemoryID(ctx, id)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "n.png", f.OriginalName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteFile_TableDriven(t *testing.T) {
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
				m.ExpectExec("UPDATE .*memory_files.*SET .*deleted_at.*").WithArgs(AnyTime{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "delete error",
			prepare: func(m sqlmock.Sqlmock, id uuid.UUID) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE .*memory_files.*SET .*deleted_at.*").WithArgs(AnyTime{}, id).WillReturnError(errors.New("delete failed"))
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
			err := r.DeleteFile(ctx, id)
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
			got, err := r.CreateFile(ctx, f)
			require.NoError(t, err)
			require.Equal(t, "tx.png", got.OriginalName)
			require.NoError(t, mockTx.ExpectationsWereMet())
			require.NoError(t, mockMain.ExpectationsWereMet())
		})
	}
}
