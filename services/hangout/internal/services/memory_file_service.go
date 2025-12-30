package services

import (
	"context"
	"path/filepath"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/storage"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryFileService interface {
	UploadFile(ctx context.Context, tx *gorm.DB, fileData *dto.FileUploadData) (*domain.MemoryFile, error)
	GetFileByMemoryID(ctx context.Context, memoryID uuid.UUID) (*dto.MemoryFileResponse, error)
	DeleteFile(ctx context.Context, tx *gorm.DB, memoryID uuid.UUID) error
	GeneratePresignedURL(ctx context.Context, storagePath string) (string, error)
}

type memoryFileService struct {
	storage       storage.Storage
	fileRepo      repository.MemoryFileRepository
	fileValidator *validator.FileValidator
}

func NewMemoryFileService(storage storage.Storage, fileRepo repository.MemoryFileRepository, fileValidator *validator.FileValidator) MemoryFileService {
	return &memoryFileService{
		storage:       storage,
		fileRepo:      fileRepo,
		fileValidator: fileValidator,
	}
}

func (s *memoryFileService) UploadFile(ctx context.Context, tx *gorm.DB, fileData *dto.FileUploadData) (*domain.MemoryFile, error) {
	if err := s.fileValidator.ValidateFileMetadata(fileData.Filename, fileData.Size, fileData.MimeType); err != nil {
		return nil, err
	}

	if err := s.storage.Upload(ctx, fileData.StoragePath, fileData.Content, fileData.MimeType); err != nil {
		return nil, err
	}

	fileExtension := filepath.Ext(fileData.Filename)
	memoryFile := &domain.MemoryFile{
		OriginalName:  fileData.Filename,
		FileExtension: fileExtension,
		StoragePath:   fileData.StoragePath,
		FileSize:      fileData.Size,
		MimeType:      fileData.MimeType,
		MemoryID:      fileData.MemoryID,
	}

	createdFile, err := s.fileRepo.WithTx(tx).CreateFile(ctx, memoryFile)
	if err != nil {
		_ = s.storage.Delete(ctx, fileData.StoragePath)
		return nil, err
	}

	return createdFile, nil
}

func (s *memoryFileService) GetFileByMemoryID(ctx context.Context, memoryID uuid.UUID) (*dto.MemoryFileResponse, error) {
	memoryFile, err := s.fileRepo.GetFileByMemoryID(ctx, memoryID)
	if err != nil {
		return nil, err
	}

	fileURL, err := s.storage.GeneratePresignedURL(ctx, memoryFile.StoragePath)
	if err != nil {
		return nil, err
	}

	return &dto.MemoryFileResponse{
		ID:            memoryFile.ID,
		MemoryID:      memoryFile.MemoryID,
		OriginalName:  memoryFile.OriginalName,
		FileExtension: memoryFile.FileExtension,
		FileSize:      memoryFile.FileSize,
		MimeType:      memoryFile.MimeType,
		FileURL:       fileURL,
	}, nil
}

func (s *memoryFileService) DeleteFile(ctx context.Context, tx *gorm.DB, memoryID uuid.UUID) error {
	memoryFile, err := s.fileRepo.GetFileByMemoryID(ctx, memoryID)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, memoryFile.StoragePath); err != nil {
		return err
	}

	if err := s.fileRepo.WithTx(tx).DeleteFile(ctx, memoryID); err != nil {
		return err
	}

	return nil
}

func (s *memoryFileService) GeneratePresignedURL(ctx context.Context, storagePath string) (string, error) {
	return s.storage.GeneratePresignedURL(ctx, storagePath)
}
