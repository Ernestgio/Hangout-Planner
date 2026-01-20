package services

import (
	"context"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/grpc"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryServiceV2 interface {
	GenerateUploadURLs(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, req *dto.GenerateUploadURLsRequest) (*dto.MemoryUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID uuid.UUID, req *dto.ConfirmUploadRequest) error
	GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error)
	ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error)
	DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error
}

type memoryServiceV2 struct {
	db          *gorm.DB
	memoryRepo  repository.MemoryRepository
	hangoutRepo repository.HangoutRepository
	fileService grpc.FileService
}

func NewMemoryServiceV2(db *gorm.DB, memoryRepo repository.MemoryRepository, hangoutRepo repository.HangoutRepository, fileService grpc.FileService,
) MemoryServiceV2 {
	return &memoryServiceV2{
		db:          db,
		memoryRepo:  memoryRepo,
		hangoutRepo: hangoutRepo,
		fileService: fileService,
	}
}

func (s *memoryServiceV2) GenerateUploadURLs(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, req *dto.GenerateUploadURLsRequest) (*dto.MemoryUploadResponse, error) {
	if len(req.Files) > constants.MaxFilePerUpload {
		return nil, apperrors.ErrTooManyFiles
	}

	_, err := s.hangoutRepo.GetHangoutByID(ctx, hangoutID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrInvalidHangoutID
		}
		return nil, err
	}

	memories := make([]*domain.Memory, len(req.Files))

	for i, file := range req.Files {
		memories[i] = &domain.Memory{
			Name:      file.Filename,
			HangoutID: hangoutID,
			UserID:    userID,
		}
	}

	var uploadURLsResp *filepb.GenerateUploadURLsResponse

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.memoryRepo.WithTx(tx).CreateMemoriesBatch(ctx, memories); err != nil {
			return err
		}

		fileIntents := make([]*filepb.FileUploadIntent, len(req.Files))
		for i, memory := range memories {
			fileIntents[i] = &filepb.FileUploadIntent{
				Filename: req.Files[i].Filename,
				Size:     req.Files[i].Size,
				MimeType: req.Files[i].MimeType,
				MemoryId: memory.ID.String(),
			}
		}

		baseStoragePath := "hangouts/" + hangoutID.String() + "/memories"
		var err error
		uploadURLsResp, err = s.fileService.GenerateUploadURLs(ctx, baseStoragePath, fileIntents)
		if err != nil {
			return err
		}

		fileIDUpdates := make(map[uuid.UUID]uuid.UUID)
		for _, presignedURL := range uploadURLsResp.Urls {
			memoryID, _ := uuid.Parse(presignedURL.MemoryId)
			fileID, _ := uuid.Parse(presignedURL.FileId)
			fileIDUpdates[memoryID] = fileID
		}

		if err := s.memoryRepo.WithTx(tx).UpdateFileIDs(ctx, fileIDUpdates); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToMemoryUploadResponse(uploadURLsResp.Urls), nil
}

func (s *memoryServiceV2) ConfirmUpload(ctx context.Context, userID uuid.UUID, req *dto.ConfirmUploadRequest) error {
	memories, err := s.memoryRepo.GetMemoriesByIDs(ctx, req.MemoryIDs, userID)
	if err != nil {
		return err
	}

	if len(memories) != len(req.MemoryIDs) {
		return apperrors.ErrMemoryNotFound
	}

	fileIDs := make([]string, 0, len(memories))
	for _, memory := range memories {
		if memory.FileID == nil {
			return apperrors.ErrMemoryNotFound
		}
		fileIDs = append(fileIDs, memory.FileID.String())
	}

	return s.fileService.ConfirmUpload(ctx, fileIDs)
}

func (s *memoryServiceV2) GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error) {
	memory, err := s.memoryRepo.GetMemoryByID(ctx, memoryID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrMemoryNotFound
		}
		return nil, err
	}

	fileWithURL, err := s.fileService.GetFileByMemoryID(ctx, memoryID.String())
	if err != nil {
		return nil, err
	}

	return mapper.MemoryToResponseDTO(memory, fileWithURL.DownloadUrl, fileWithURL.FileSize, fileWithURL.MimeType), nil
}

func (s *memoryServiceV2) ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error) {
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

	memoryIDs := make([]string, len(memories))
	for i, memory := range memories {
		memoryIDs[i] = memory.ID.String()
	}

	filesMap, err := s.fileService.GetFilesByMemoryIDs(ctx, memoryIDs)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MemoryResponse, 0, len(memories))
	for _, memory := range memories {
		fileWithURL := filesMap[memory.ID.String()]
		if fileWithURL != nil {
			responses = append(responses, *mapper.MemoryToResponseDTO(
				&memory,
				fileWithURL.DownloadUrl,
				fileWithURL.FileSize,
				fileWithURL.MimeType,
			))
		}
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

func (s *memoryServiceV2) DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.memoryRepo.WithTx(tx).GetMemoryByID(ctx, memoryID, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return apperrors.ErrMemoryNotFound
			}
			return err
		}

		if err := s.fileService.DeleteFile(ctx, memoryID.String()); err != nil {
			return err
		}

		if err := s.memoryRepo.WithTx(tx).DeleteMemory(ctx, memoryID); err != nil {
			return err
		}

		return nil
	})
}
