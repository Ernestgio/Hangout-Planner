package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMemoryService_GenerateUploadURLs(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	hangoutID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *dto.GenerateUploadURLsRequest
		setup     func(*MockMemoryRepository, *MockHangoutRepository, *MockFileService, sqlmock.Sqlmock)
		wantError error
	}{
		{
			name: "success",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("CreateMemoriesBatch", ctx, mock.Anything).Return(nil)
				fileService.On("GenerateUploadURLs", ctx, "hangouts/"+hangoutID.String()+"/memories", mock.Anything).Return(&filepb.GenerateUploadURLsResponse{
					Urls: []*filepb.PresignedUploadURL{
						{FileId: uuid.New().String(), MemoryId: uuid.New().String(), Filename: "photo.jpg", UploadUrl: "https://s3/upload", ExpiresAt: 123456789},
					},
				}, nil)
				memRepo.On("UpdateFileIDs", ctx, mock.Anything).Return(nil)
				sqlMock.ExpectCommit()
			},
		},
		{
			name: "too many files",
			req: &dto.GenerateUploadURLsRequest{
				Files: make([]dto.FileUploadIntent, 11),
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
			},
			wantError: apperrors.ErrTooManyFiles,
		},
		{
			name: "hangout not found",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrInvalidHangoutID,
		},
		{
			name: "hangout repo error",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, dbError)
			},
			wantError: dbError,
		},
		{
			name: "create memories batch error",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("CreateMemoriesBatch", ctx, mock.Anything).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
		{
			name: "file service error",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("CreateMemoriesBatch", ctx, mock.Anything).Return(nil)
				fileService.On("GenerateUploadURLs", ctx, "hangouts/"+hangoutID.String()+"/memories", mock.Anything).Return(nil, dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
		{
			name: "update file ids error",
			req: &dto.GenerateUploadURLsRequest{
				Files: []dto.FileUploadIntent{
					{Filename: "photo.jpg", Size: 1024, MimeType: "image/jpeg"},
				},
			},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("CreateMemoriesBatch", ctx, mock.Anything).Return(nil)
				fileService.On("GenerateUploadURLs", ctx, "hangouts/"+hangoutID.String()+"/memories", mock.Anything).Return(&filepb.GenerateUploadURLsResponse{
					Urls: []*filepb.PresignedUploadURL{
						{FileId: uuid.New().String(), MemoryId: uuid.New().String(), Filename: "photo.jpg", UploadUrl: "https://s3/upload", ExpiresAt: 123456789},
					},
				}, nil)
				memRepo.On("UpdateFileIDs", ctx, mock.Anything).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			memRepo := new(MockMemoryRepository)
			hangoutRepo := new(MockHangoutRepository)
			fileService := new(MockFileService)
			tt.setup(memRepo, hangoutRepo, fileService, sqlMock)
			svc := services.NewMemoryService(db, memRepo, hangoutRepo, fileService, nil)
			resp, err := svc.GenerateUploadURLs(ctx, userID, hangoutID, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
			memRepo.AssertExpectations(t)
			hangoutRepo.AssertExpectations(t)
			fileService.AssertExpectations(t)
		})
	}
}

func TestMemoryService_ConfirmUpload(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	memoryID1 := uuid.New()
	memoryID2 := uuid.New()
	fileID1 := uuid.New()
	fileID2 := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		req       *dto.ConfirmUploadRequest
		setup     func(*MockMemoryRepository, *MockFileService)
		wantError error
	}{
		{
			name: "success",
			req: &dto.ConfirmUploadRequest{
				MemoryIDs: []uuid.UUID{memoryID1, memoryID2},
			},
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoriesByIDs", ctx, []uuid.UUID{memoryID1, memoryID2}, userID).Return([]domain.Memory{
					{ID: memoryID1, FileID: &fileID1},
					{ID: memoryID2, FileID: &fileID2},
				}, nil)
				fileService.On("ConfirmUpload", ctx, []string{fileID1.String(), fileID2.String()}).Return(nil)
			},
		},
		{
			name: "get memories error",
			req: &dto.ConfirmUploadRequest{
				MemoryIDs: []uuid.UUID{memoryID1},
			},
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoriesByIDs", ctx, []uuid.UUID{memoryID1}, userID).Return(nil, dbError)
			},
			wantError: dbError,
		},
		{
			name: "memories count mismatch",
			req: &dto.ConfirmUploadRequest{
				MemoryIDs: []uuid.UUID{memoryID1, memoryID2},
			},
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoriesByIDs", ctx, []uuid.UUID{memoryID1, memoryID2}, userID).Return([]domain.Memory{
					{ID: memoryID1, FileID: &fileID1},
				}, nil)
			},
			wantError: apperrors.ErrMemoryNotFound,
		},
		{
			name: "missing file id",
			req: &dto.ConfirmUploadRequest{
				MemoryIDs: []uuid.UUID{memoryID1},
			},
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoriesByIDs", ctx, []uuid.UUID{memoryID1}, userID).Return([]domain.Memory{
					{ID: memoryID1, FileID: nil},
				}, nil)
			},
			wantError: apperrors.ErrMemoryNotFound,
		},
		{
			name: "file service error",
			req: &dto.ConfirmUploadRequest{
				MemoryIDs: []uuid.UUID{memoryID1},
			},
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoriesByIDs", ctx, []uuid.UUID{memoryID1}, userID).Return([]domain.Memory{
					{ID: memoryID1, FileID: &fileID1},
				}, nil)
				fileService.On("ConfirmUpload", ctx, []string{fileID1.String()}).Return(dbError)
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupDB(t)
			memRepo := new(MockMemoryRepository)
			fileService := new(MockFileService)
			tt.setup(memRepo, fileService)
			svc := services.NewMemoryService(db, memRepo, nil, fileService, nil)
			err := svc.ConfirmUpload(ctx, userID, tt.req)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}
			memRepo.AssertExpectations(t)
			fileService.AssertExpectations(t)
		})
	}
}

func TestMemoryService_GetMemory(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	memoryID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		setup     func(*MockMemoryRepository, *MockFileService)
		wantError error
	}{
		{
			name: "success",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(&domain.Memory{ID: memoryID, Name: "photo.jpg"}, nil)
				fileService.On("GetFileByMemoryID", ctx, memoryID.String()).Return(&filepb.FileWithURL{
					DownloadUrl: "https://s3/download",
					FileSize:    1024,
					MimeType:    "image/jpeg",
				}, nil)
			},
		},
		{
			name: "memory not found",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrMemoryNotFound,
		},
		{
			name: "memory repo error",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(nil, dbError)
			},
			wantError: dbError,
		},
		{
			name: "file service error",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService) {
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(&domain.Memory{ID: memoryID, Name: "photo.jpg"}, nil)
				fileService.On("GetFileByMemoryID", ctx, memoryID.String()).Return(nil, dbError)
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupDB(t)
			memRepo := new(MockMemoryRepository)
			fileService := new(MockFileService)
			tt.setup(memRepo, fileService)
			svc := services.NewMemoryService(db, memRepo, nil, fileService, nil)
			resp, err := svc.GetMemory(ctx, userID, memoryID)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
			memRepo.AssertExpectations(t)
			fileService.AssertExpectations(t)
		})
	}
}

func TestMemoryService_ListMemories(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	hangoutID := uuid.New()
	memoryID1 := uuid.New()
	memoryID2 := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name       string
		pagination *dto.CursorPagination
		setup      func(*MockMemoryRepository, *MockHangoutRepository, *MockFileService)
		wantError  error
		wantMore   bool
	}{
		{
			name:       "success with results",
			pagination: &dto.CursorPagination{Limit: 2},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				memRepo.On("GetMemoriesByHangoutID", ctx, hangoutID, mock.Anything).Return([]domain.Memory{
					{ID: memoryID1, Name: "photo1.jpg"},
					{ID: memoryID2, Name: "photo2.jpg"},
				}, nil)
				fileService.On("GetFilesByMemoryIDs", ctx, []string{memoryID1.String(), memoryID2.String()}).Return(map[string]*filepb.FileWithURL{
					memoryID1.String(): {DownloadUrl: "https://s3/file1", FileSize: 1024, MimeType: "image/jpeg"},
					memoryID2.String(): {DownloadUrl: "https://s3/file2", FileSize: 2048, MimeType: "image/png"},
				}, nil)
			},
		},
		{
			name:       "success with more results",
			pagination: &dto.CursorPagination{Limit: 1},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				memRepo.On("GetMemoriesByHangoutID", ctx, hangoutID, mock.Anything).Return([]domain.Memory{
					{ID: memoryID1, Name: "photo1.jpg"},
					{ID: memoryID2, Name: "photo2.jpg"},
				}, nil)
				fileService.On("GetFilesByMemoryIDs", ctx, []string{memoryID1.String()}).Return(map[string]*filepb.FileWithURL{
					memoryID1.String(): {DownloadUrl: "https://s3/file1", FileSize: 1024, MimeType: "image/jpeg"},
				}, nil)
			},
			wantMore: true,
		},
		{
			name:       "hangout not found",
			pagination: &dto.CursorPagination{Limit: 2},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, gorm.ErrRecordNotFound)
			},
			wantError: apperrors.ErrInvalidHangoutID,
		},
		{
			name:       "hangout repo error",
			pagination: &dto.CursorPagination{Limit: 2},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(nil, dbError)
			},
			wantError: dbError,
		},
		{
			name:       "memory repo error",
			pagination: &dto.CursorPagination{Limit: 2},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				memRepo.On("GetMemoriesByHangoutID", ctx, hangoutID, mock.Anything).Return(nil, dbError)
			},
			wantError: dbError,
		},
		{
			name:       "file service error",
			pagination: &dto.CursorPagination{Limit: 2},
			setup: func(memRepo *MockMemoryRepository, hangoutRepo *MockHangoutRepository, fileService *MockFileService) {
				hangoutRepo.On("GetHangoutByID", ctx, hangoutID, userID).Return(&domain.Hangout{ID: hangoutID}, nil)
				memRepo.On("GetMemoriesByHangoutID", ctx, hangoutID, mock.Anything).Return([]domain.Memory{
					{ID: memoryID1, Name: "photo1.jpg"},
				}, nil)
				fileService.On("GetFilesByMemoryIDs", ctx, []string{memoryID1.String()}).Return(nil, dbError)
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupDB(t)
			memRepo := new(MockMemoryRepository)
			hangoutRepo := new(MockHangoutRepository)
			fileService := new(MockFileService)
			tt.setup(memRepo, hangoutRepo, fileService)
			svc := services.NewMemoryService(db, memRepo, hangoutRepo, fileService, nil)
			resp, err := svc.ListMemories(ctx, userID, hangoutID, tt.pagination)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.wantMore, resp.HasMore)
			}
			memRepo.AssertExpectations(t)
			hangoutRepo.AssertExpectations(t)
			fileService.AssertExpectations(t)
		})
	}
}

func TestMemoryService_DeleteMemory(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	memoryID := uuid.New()
	dbError := errors.New("db error")

	tests := []struct {
		name      string
		setup     func(*MockMemoryRepository, *MockFileService, sqlmock.Sqlmock)
		wantError error
	}{
		{
			name: "success",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(&domain.Memory{ID: memoryID}, nil)
				fileService.On("DeleteFile", ctx, memoryID.String()).Return(nil)
				memRepo.On("DeleteMemory", ctx, memoryID).Return(nil)
				sqlMock.ExpectCommit()
			},
		},
		{
			name: "memory not found",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(nil, gorm.ErrRecordNotFound)
				sqlMock.ExpectRollback()
			},
			wantError: apperrors.ErrMemoryNotFound,
		},
		{
			name: "get memory error",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(nil, dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
		{
			name: "file service error",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(&domain.Memory{ID: memoryID}, nil)
				fileService.On("DeleteFile", ctx, memoryID.String()).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
		{
			name: "delete memory error",
			setup: func(memRepo *MockMemoryRepository, fileService *MockFileService, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				memRepo.On("WithTx", mock.Anything).Return(memRepo)
				memRepo.On("GetMemoryByID", ctx, memoryID, userID).Return(&domain.Memory{ID: memoryID}, nil)
				fileService.On("DeleteFile", ctx, memoryID.String()).Return(nil)
				memRepo.On("DeleteMemory", ctx, memoryID).Return(dbError)
				sqlMock.ExpectRollback()
			},
			wantError: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := setupDB(t)
			memRepo := new(MockMemoryRepository)
			fileService := new(MockFileService)
			tt.setup(memRepo, fileService, sqlMock)
			svc := services.NewMemoryService(db, memRepo, nil, fileService, nil)
			err := svc.DeleteMemory(ctx, userID, memoryID)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}
			memRepo.AssertExpectations(t)
			fileService.AssertExpectations(t)
		})
	}
}
