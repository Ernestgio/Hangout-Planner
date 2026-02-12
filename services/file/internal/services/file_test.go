package services_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/services"
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

type MockMemoryFileRepository struct {
	mock.Mock
}

func (m *MockMemoryFileRepository) WithTx(tx *gorm.DB) repository.MemoryFileRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.MemoryFileRepository)
}

func (m *MockMemoryFileRepository) Create(ctx context.Context, file *domain.MemoryFile) (*domain.MemoryFile, error) {
	args := m.Called(ctx, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MemoryFile), args.Error(1)
}

func (m *MockMemoryFileRepository) CreateBatch(ctx context.Context, files []*domain.MemoryFile) error {
	args := m.Called(ctx, files)
	return args.Error(0)
}

func (m *MockMemoryFileRepository) GetByID(ctx context.Context, fileID uuid.UUID) (*domain.MemoryFile, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MemoryFile), args.Error(1)
}

func (m *MockMemoryFileRepository) GetByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error) {
	args := m.Called(ctx, memoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MemoryFile), args.Error(1)
}

func (m *MockMemoryFileRepository) GetByMemoryIDs(ctx context.Context, memoryIDs []uuid.UUID) ([]*domain.MemoryFile, error) {
	args := m.Called(ctx, memoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.MemoryFile), args.Error(1)
}

func (m *MockMemoryFileRepository) UpdateStatusBatch(ctx context.Context, fileIDs []uuid.UUID, status string) error {
	args := m.Called(ctx, fileIDs, status)
	return args.Error(0)
}

func (m *MockMemoryFileRepository) Delete(ctx context.Context, memoryID uuid.UUID) error {
	args := m.Called(ctx, memoryID)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, path string, reader io.Reader, contentType string) error {
	args := m.Called(ctx, path, reader, contentType)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockStorage) GeneratePresignedDownloadURL(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GeneratePresignedUploadURL(ctx context.Context, path string, contentType string) (string, error) {
	args := m.Called(ctx, path, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetPresignedURLExpiry() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

type MockFileValidator struct {
	mock.Mock
}

func (m *MockFileValidator) ValidateFileUploadIntent(filename string, size int64, mimeType string) error {
	args := m.Called(filename, size, mimeType)
	return args.Error(0)
}

func (m *MockFileValidator) GetMaxFileSize() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockFileValidator) IsExtensionAllowed(extension string) bool {
	args := m.Called(extension)
	return args.Bool(0)
}

func TestFileService_GenerateUploadURLs(t *testing.T) {
	ctx := context.Background()
	memoryID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *filepb.GenerateUploadURLsRequest
		setup     func(*MockMemoryFileRepository, *MockStorage, *MockFileValidator, sqlmock.Sqlmock)
		wantError error
	}{
		{
			name: "success",
			req: &filepb.GenerateUploadURLsRequest{
				BaseStoragePath: "hangouts/123/memories",
				Files: []*filepb.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg", MemoryId: memoryID.String()},
				},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, val *MockFileValidator, sqlMock sqlmock.Sqlmock) {
				val.On("ValidateFileUploadIntent", "photo.jpg", int64(1024), "image/jpeg").Return(nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("CreateBatch", ctx, mock.Anything).Return(nil)
				store.On("GetPresignedURLExpiry").Return(1 * time.Hour)
				store.On("GeneratePresignedUploadURL", ctx, mock.Anything, "image/jpeg").Return("https://s3/upload", nil)
				sqlMock.ExpectCommit()
			},
		},
		{
			name: "validation error",
			req: &filepb.GenerateUploadURLsRequest{
				BaseStoragePath: "hangouts/123/memories",
				Files: []*filepb.FileUploadIntent{
					{Filename: "photo.jpg", Size: 0, MimeType: "image/jpeg", MemoryId: memoryID.String()},
				},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, val *MockFileValidator, sqlMock sqlmock.Sqlmock) {
				val.On("ValidateFileUploadIntent", "photo.jpg", int64(0), "image/jpeg").Return(apperrors.ErrInvalidFileSize)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrInvalidFileSize,
		},
		{
			name: "invalid memory id",
			req: &filepb.GenerateUploadURLsRequest{
				BaseStoragePath: "hangouts/123/memories",
				Files: []*filepb.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg", MemoryId: "invalid-uuid"},
				},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, val *MockFileValidator, sqlMock sqlmock.Sqlmock) {
				val.On("ValidateFileUploadIntent", "photo.jpg", int64(1024), "image/jpeg").Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "create batch error",
			req: &filepb.GenerateUploadURLsRequest{
				BaseStoragePath: "hangouts/123/memories",
				Files: []*filepb.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg", MemoryId: memoryID.String()},
				},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, val *MockFileValidator, sqlMock sqlmock.Sqlmock) {
				val.On("ValidateFileUploadIntent", "photo.jpg", int64(1024), "image/jpeg").Return(nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("CreateBatch", ctx, mock.Anything).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrFileCreationFailed,
		},
		{
			name: "presigned url error",
			req: &filepb.GenerateUploadURLsRequest{
				BaseStoragePath: "hangouts/123/memories",
				Files: []*filepb.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg", MemoryId: memoryID.String()},
				},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, val *MockFileValidator, sqlMock sqlmock.Sqlmock) {
				val.On("ValidateFileUploadIntent", "photo.jpg", int64(1024), "image/jpeg").Return(nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("CreateBatch", ctx, mock.Anything).Return(nil)
				store.On("GetPresignedURLExpiry").Return(1 * time.Hour)
				store.On("GeneratePresignedUploadURL", ctx, mock.Anything, "image/jpeg").Return("", apperrors.ErrPresignedUploadURLFailed)
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrPresignedUploadURLFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			repo := new(MockMemoryFileRepository)
			store := new(MockStorage)
			val := new(MockFileValidator)
			tt.setup(repo, store, val, sqlMock)
			svc := services.NewFileService(db, repo, store, val, nil)
			resp, err := svc.GenerateUploadURLs(ctx, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Len(t, resp.Urls, len(tt.req.Files))
			}
			repo.AssertExpectations(t)
			store.AssertExpectations(t)
			val.AssertExpectations(t)
		})
	}
}

func TestFileService_ConfirmUpload(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *filepb.ConfirmUploadRequest
		setup     func(*MockMemoryFileRepository, sqlmock.Sqlmock)
		wantError error
	}{
		{
			name: "success",
			req: &filepb.ConfirmUploadRequest{
				FileIds: []string{fileID.String()},
			},
			setup: func(repo *MockMemoryFileRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("UpdateStatusBatch", ctx, []uuid.UUID{fileID}, string(enums.FileUploadStatusUploaded)).Return(nil)
				sqlMock.ExpectCommit()
			},
		},
		{
			name: "invalid uuid",
			req: &filepb.ConfirmUploadRequest{
				FileIds: []string{"invalid-uuid"},
			},
			setup:     func(repo *MockMemoryFileRepository, sqlMock sqlmock.Sqlmock) {},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "update status error",
			req: &filepb.ConfirmUploadRequest{
				FileIds: []string{fileID.String()},
			},
			setup: func(repo *MockMemoryFileRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("UpdateStatusBatch", ctx, []uuid.UUID{fileID}, string(enums.FileUploadStatusUploaded)).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrFileStatusUpdateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			repo := new(MockMemoryFileRepository)
			tt.setup(repo, sqlMock)
			svc := services.NewFileService(db, repo, nil, nil, nil)
			resp, err := svc.ConfirmUpload(ctx, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Success)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestFileService_GetFileByMemoryID(t *testing.T) {
	ctx := context.Background()
	memoryID := uuid.New()
	fileID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *filepb.GetFileByMemoryIDRequest
		setup     func(*MockMemoryFileRepository, *MockStorage)
		wantError error
	}{
		{
			name: "success",
			req: &filepb.GetFileByMemoryIDRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(&domain.MemoryFile{
					ID:           fileID,
					MemoryID:     memoryID,
					OriginalName: "photo.jpg",
					StoragePath:  "path/photo.jpg",
					FileSize:     1024,
					MimeType:     "image/jpeg",
				}, nil)
				store.On("GeneratePresignedDownloadURL", ctx, "path/photo.jpg").Return("https://s3/download", nil)
				store.On("GetPresignedURLExpiry").Return(1 * time.Hour)
			},
		},
		{
			name: "invalid uuid",
			req: &filepb.GetFileByMemoryIDRequest{
				MemoryId: "invalid-uuid",
			},
			setup:     func(repo *MockMemoryFileRepository, store *MockStorage) {},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "file not found",
			req: &filepb.GetFileByMemoryIDRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "presigned url error",
			req: &filepb.GetFileByMemoryIDRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(&domain.MemoryFile{
					ID:           fileID,
					MemoryID:     memoryID,
					StoragePath:  "path/photo.jpg",
					OriginalName: "photo.jpg",
				}, nil)
				store.On("GeneratePresignedDownloadURL", ctx, "path/photo.jpg").Return("", dbError)
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupDB(t)
			repo := new(MockMemoryFileRepository)
			store := new(MockStorage)
			tt.setup(repo, store)
			svc := services.NewFileService(db, repo, store, nil, nil)
			resp, err := svc.GetFileByMemoryID(ctx, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.File)
			}
			repo.AssertExpectations(t)
			store.AssertExpectations(t)
		})
	}
}

func TestFileService_GetFilesByMemoryIDs(t *testing.T) {
	ctx := context.Background()
	memoryID1 := uuid.New()
	memoryID2 := uuid.New()
	fileID1 := uuid.New()
	fileID2 := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *filepb.GetFilesByMemoryIDsRequest
		setup     func(*MockMemoryFileRepository, *MockStorage)
		wantError error
	}{
		{
			name: "success",
			req: &filepb.GetFilesByMemoryIDsRequest{
				MemoryIds: []string{memoryID1.String(), memoryID2.String()},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryIDs", ctx, []uuid.UUID{memoryID1, memoryID2}).Return([]*domain.MemoryFile{
					{ID: fileID1, MemoryID: memoryID1, StoragePath: "path1/photo1.jpg", OriginalName: "photo1.jpg"},
					{ID: fileID2, MemoryID: memoryID2, StoragePath: "path2/photo2.jpg", OriginalName: "photo2.jpg"},
				}, nil)
				store.On("GeneratePresignedDownloadURL", ctx, "path1/photo1.jpg").Return("https://s3/download1", nil)
				store.On("GeneratePresignedDownloadURL", ctx, "path2/photo2.jpg").Return("https://s3/download2", nil)
				store.On("GetPresignedURLExpiry").Return(1 * time.Hour)
			},
		},
		{
			name: "invalid uuid",
			req: &filepb.GetFilesByMemoryIDsRequest{
				MemoryIds: []string{"invalid-uuid"},
			},
			setup:     func(repo *MockMemoryFileRepository, store *MockStorage) {},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "files not found",
			req: &filepb.GetFilesByMemoryIDsRequest{
				MemoryIds: []string{memoryID1.String()},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryIDs", ctx, []uuid.UUID{memoryID1}).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "presigned url error",
			req: &filepb.GetFilesByMemoryIDsRequest{
				MemoryIds: []string{memoryID1.String()},
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage) {
				repo.On("GetByMemoryIDs", ctx, []uuid.UUID{memoryID1}).Return([]*domain.MemoryFile{
					{ID: fileID1, MemoryID: memoryID1, StoragePath: "path1/photo1.jpg"},
				}, nil)
				store.On("GeneratePresignedDownloadURL", ctx, "path1/photo1.jpg").Return("", dbError)
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupDB(t)
			repo := new(MockMemoryFileRepository)
			store := new(MockStorage)
			tt.setup(repo, store)
			svc := services.NewFileService(db, repo, store, nil, nil)
			resp, err := svc.GetFilesByMemoryIDs(ctx, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.Files)
			}
			repo.AssertExpectations(t)
			store.AssertExpectations(t)
		})
	}
}

func TestFileService_DeleteFile(t *testing.T) {
	ctx := context.Background()
	memoryID := uuid.New()
	fileID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *filepb.DeleteFileRequest
		setup     func(*MockMemoryFileRepository, *MockStorage, sqlmock.Sqlmock)
		wantError error
	}{
		{
			name: "success",
			req: &filepb.DeleteFileRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, sqlMock sqlmock.Sqlmock) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(&domain.MemoryFile{
					ID:          fileID,
					MemoryID:    memoryID,
					StoragePath: "path/photo.jpg",
				}, nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("Delete", ctx, fileID).Return(nil)
				sqlMock.ExpectCommit()
				store.On("Delete", ctx, "path/photo.jpg").Return(nil)
			},
		},
		{
			name: "invalid uuid",
			req: &filepb.DeleteFileRequest{
				MemoryId: "invalid-uuid",
			},
			setup:     func(repo *MockMemoryFileRepository, store *MockStorage, sqlMock sqlmock.Sqlmock) {},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "file not found",
			req: &filepb.DeleteFileRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, sqlMock sqlmock.Sqlmock) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrInvalidMemoryID,
		},
		{
			name: "delete from db error",
			req: &filepb.DeleteFileRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, sqlMock sqlmock.Sqlmock) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(&domain.MemoryFile{
					ID:          fileID,
					MemoryID:    memoryID,
					StoragePath: "path/photo.jpg",
				}, nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("Delete", ctx, fileID).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrFileDeleteFailed,
		},
		{
			name: "delete from storage error",
			req: &filepb.DeleteFileRequest{
				MemoryId: memoryID.String(),
			},
			setup: func(repo *MockMemoryFileRepository, store *MockStorage, sqlMock sqlmock.Sqlmock) {
				repo.On("GetByMemoryID", ctx, memoryID).Return(&domain.MemoryFile{
					ID:          fileID,
					MemoryID:    memoryID,
					StoragePath: "path/photo.jpg",
				}, nil)
				sqlMock.ExpectBegin()
				repo.On("WithTx", mock.Anything).Return(repo)
				repo.On("Delete", ctx, fileID).Return(nil)
				sqlMock.ExpectCommit()
				store.On("Delete", ctx, "path/photo.jpg").Return(dbError)
			},
			wantError: apperrors.ErrFileDeleteFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			repo := new(MockMemoryFileRepository)
			store := new(MockStorage)
			tt.setup(repo, store, sqlMock)
			svc := services.NewFileService(db, repo, store, nil, nil)
			resp, err := svc.DeleteFile(ctx, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Success)
			}
			repo.AssertExpectations(t)
			store.AssertExpectations(t)
		})
	}
}
