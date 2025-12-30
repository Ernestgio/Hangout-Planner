package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryService interface {
	CreateMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, files []*multipart.FileHeader) ([]dto.MemoryResponse, error)
	GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error)
	ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error)
	DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error
}

type memoryService struct {
	db                *gorm.DB
	memoryRepo        repository.MemoryRepository
	hangoutRepo       repository.HangoutRepository
	memoryFileService MemoryFileService
}

func NewMemoryService(db *gorm.DB, memoryRepo repository.MemoryRepository, hangoutRepo repository.HangoutRepository, memoryFileService MemoryFileService) MemoryService {
	return &memoryService{
		db:                db,
		memoryRepo:        memoryRepo,
		hangoutRepo:       hangoutRepo,
		memoryFileService: memoryFileService,
	}
}

func (s *memoryService) CreateMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, files []*multipart.FileHeader) ([]dto.MemoryResponse, error) {
	if len(files) > constants.MaxFilePerUpload {
		return nil, apperrors.ErrTooManyFiles
	}

	_, err := s.hangoutRepo.GetHangoutByID(ctx, hangoutID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrInvalidHangoutID
		}
		return nil, err
	}

	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		responses   []dto.MemoryResponse
		failedCount int
	)

	for _, fileHeader := range files {
		wg.Add(1)

		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				filename, size, mimeType := validator.ExtractFileMetadata(fh)

				memory := &domain.Memory{
					Name:      filename,
					HangoutID: hangoutID,
					UserID:    userID,
				}

				createdMemory, err := s.memoryRepo.WithTx(tx).CreateMemory(ctx, memory)
				if err != nil {
					return err
				}

				src, err := fh.Open()
				if err != nil {
					return err
				}
				defer func() {
					_ = src.Close()
				}()

				storagePath := fmt.Sprintf("hangouts/%s/memories/%s/%s", hangoutID, createdMemory.ID, filename)

				memoryFile, err := s.memoryFileService.UploadFile(ctx, tx, &dto.FileUploadData{
					MemoryID:    createdMemory.ID,
					Filename:    filename,
					StoragePath: storagePath,
					Size:        size,
					MimeType:    mimeType,
					Content:     src,
				})
				if err != nil {
					return err
				}

				fileURL, err := s.memoryFileService.GeneratePresignedURL(ctx, memoryFile.StoragePath)
				if err != nil {
					return err
				}

				mu.Lock()
				responses = append(responses, *mapper.MemoryToResponseDTO(
					createdMemory,
					fileURL,
					memoryFile.FileSize,
					memoryFile.MimeType,
				))
				mu.Unlock()

				return nil
			})

			if err != nil {
				mu.Lock()
				failedCount++
				mu.Unlock()
			}
		}(fileHeader)
	}

	wg.Wait()

	if len(responses) > 0 {
		return responses, nil
	}

	return nil, apperrors.ErrAllFilesUploadFailed
}

func (s *memoryService) GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error) {
	memory, err := s.memoryRepo.GetMemoryByID(ctx, memoryID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrMemoryNotFound
		}
		return nil, err
	}

	memoryFile, err := s.memoryFileService.GetFileByMemoryID(ctx, memoryID)
	if err != nil {
		return nil, err
	}

	return mapper.MemoryToResponseDTO(memory, memoryFile.FileURL, memoryFile.FileSize, memoryFile.MimeType), nil
}

func (s *memoryService) ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error) {
	_, err := s.hangoutRepo.GetHangoutByID(ctx, hangoutID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrInvalidHangoutID
		}
		return nil, err
	}

	memories, err := s.memoryRepo.GetMemoriesByHangoutID(ctx, hangoutID, pagination)
	if err != nil {
		return nil, err
	}

	limit := pagination.GetLimit()
	hasMore := len(memories) > limit
	if hasMore {
		memories = memories[:limit]
	}

	responses := make([]dto.MemoryResponse, 0, len(memories))
	for _, memory := range memories {
		memoryFile, err := s.memoryFileService.GetFileByMemoryID(ctx, memory.ID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, *mapper.MemoryToResponseDTO(
			&memory,
			memoryFile.FileURL,
			memoryFile.FileSize,
			memoryFile.MimeType,
		))
	}

	var nextCursor *uuid.UUID
	if hasMore && len(memories) > 0 {
		lastID := memories[len(memories)-1].ID
		nextCursor = &lastID
	}

	return &dto.PaginatedMemories{
		Data:       responses,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *memoryService) DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.memoryRepo.WithTx(tx).GetMemoryByID(ctx, memoryID, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return apperrors.ErrMemoryNotFound
			}
			return err
		}

		if err := s.memoryFileService.DeleteFile(ctx, tx, memoryID); err != nil {
			return err
		}

		if err := s.memoryRepo.WithTx(tx).DeleteMemory(ctx, memoryID); err != nil {
			return err
		}

		return nil
	})
}
