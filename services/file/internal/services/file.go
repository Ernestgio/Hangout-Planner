package services

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService interface {
	GenerateUploadURLs(ctx context.Context, req *filepb.GenerateUploadURLsRequest) (*filepb.GenerateUploadURLsResponse, error)
	ConfirmUpload(ctx context.Context, req *filepb.ConfirmUploadRequest) (*filepb.ConfirmUploadResponse, error)
	GetFileByMemoryID(ctx context.Context, req *filepb.GetFileByMemoryIDRequest) (*filepb.GetFileByMemoryIDResponse, error)
	GetFilesByMemoryIDs(ctx context.Context, req *filepb.GetFilesByMemoryIDsRequest) (*filepb.GetFilesByMemoryIDsResponse, error)
	DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error)
}

type fileService struct {
	db       *gorm.DB
	fileRepo repository.MemoryFileRepository
	storage  storage.Storage
}

func NewFileService(db *gorm.DB, repo repository.MemoryFileRepository, storage storage.Storage) FileService {
	return &fileService{
		db:       db,
		fileRepo: repo,
		storage:  storage,
	}
}

func (s *fileService) GenerateUploadURLs(ctx context.Context, req *filepb.GenerateUploadURLsRequest) (*filepb.GenerateUploadURLsResponse, error) {
	var urls []*filepb.PresignedUploadURL

	err := s.db.Transaction(func(tx *gorm.DB) error {
		files := make([]*domain.MemoryFile, 0, len(req.Files))

		for _, intent := range req.Files {
			file, err := mapper.ToDomainMemoryFile(intent, req.BaseStoragePath, enums.FileUploadStatusPending)
			if err != nil {
				return apperrors.ErrInvalidMemoryID
			}
			files = append(files, file)
		}

		if err := s.fileRepo.WithTx(tx).CreateBatch(ctx, files); err != nil {
			return apperrors.ErrFileCreationFailed
		}

		expiresAt := mapper.GetExpiresAtUnix(s.storage.GetPresignedURLExpiry())
		urls = make([]*filepb.PresignedUploadURL, 0, len(files))

		for _, file := range files {
			uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, file.StoragePath, file.MimeType)
			if err != nil {
				return err
			}
			presignedURL := mapper.ToPresignedUploadURL(file.ID, file.OriginalName, uploadURL, expiresAt)
			urls = append(urls, presignedURL)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &filepb.GenerateUploadURLsResponse{
		Urls: urls,
	}, nil
}

func (s *fileService) ConfirmUpload(ctx context.Context, req *filepb.ConfirmUploadRequest) (*filepb.ConfirmUploadResponse, error) {
	fileIDs := make([]uuid.UUID, 0, len(req.FileIds))
	for _, idStr := range req.FileIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, apperrors.ErrInvalidMemoryID
		}
		fileIDs = append(fileIDs, id)
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.fileRepo.WithTx(tx).UpdateStatusBatch(ctx, fileIDs, string(enums.FileUploadStatusUploaded)); err != nil {
			return apperrors.ErrFileStatusUpdateFailed
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &filepb.ConfirmUploadResponse{
		Success: true,
	}, nil
}

func (s *fileService) GetFileByMemoryID(ctx context.Context, req *filepb.GetFileByMemoryIDRequest) (*filepb.GetFileByMemoryIDResponse, error) {
	memoryID, err := uuid.Parse(req.MemoryId)
	if err != nil {
		return nil, apperrors.ErrInvalidMemoryID
	}

	file, err := s.fileRepo.GetByMemoryID(ctx, memoryID)
	if err != nil {
		return nil, apperrors.ErrInvalidMemoryID
	}

	downloadURL, err := s.storage.GeneratePresignedDownloadURL(ctx, file.StoragePath)
	if err != nil {
		return nil, err
	}

	urlExpiresAt := mapper.GetExpiresAtUnix(s.storage.GetPresignedURLExpiry())
	fileWithURL := mapper.ToFileWithURL(file, downloadURL, urlExpiresAt)

	return &filepb.GetFileByMemoryIDResponse{
		File: fileWithURL,
	}, nil
}

func (s *fileService) GetFilesByMemoryIDs(ctx context.Context, req *filepb.GetFilesByMemoryIDsRequest) (*filepb.GetFilesByMemoryIDsResponse, error) {
	memoryIDs := make([]uuid.UUID, 0, len(req.MemoryIds))
	for _, idStr := range req.MemoryIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, apperrors.ErrInvalidMemoryID
		}
		memoryIDs = append(memoryIDs, id)
	}

	files, err := s.fileRepo.GetByMemoryIDs(ctx, memoryIDs)
	if err != nil {
		return nil, apperrors.ErrInvalidMemoryID
	}

	downloadURLs := make(map[uuid.UUID]string, len(files))
	for _, file := range files {
		downloadURL, err := s.storage.GeneratePresignedDownloadURL(ctx, file.StoragePath)
		if err != nil {
			return nil, err
		}
		downloadURLs[file.ID] = downloadURL
	}

	urlExpiresAt := mapper.GetExpiresAtUnix(s.storage.GetPresignedURLExpiry())
	filesWithURLs := mapper.ToFileWithURLBatch(files, downloadURLs, urlExpiresAt)

	return &filepb.GetFilesByMemoryIDsResponse{
		Files: filesWithURLs,
	}, nil
}

func (s *fileService) DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	memoryID, err := uuid.Parse(req.MemoryId)
	if err != nil {
		return nil, apperrors.ErrInvalidMemoryID
	}

	file, err := s.fileRepo.GetByMemoryID(ctx, memoryID)
	if err != nil {
		return nil, apperrors.ErrInvalidMemoryID
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.fileRepo.WithTx(tx).Delete(ctx, file.ID); err != nil {
			return apperrors.ErrFileDeleteFailed
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.storage.Delete(ctx, file.StoragePath); err != nil {
		return nil, apperrors.ErrFileDeleteFailed
	}

	return &filepb.DeleteFileResponse{
		Success: true,
	}, nil
}
