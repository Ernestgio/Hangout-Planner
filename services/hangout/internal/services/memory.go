package services

import (
	"context"
	"time"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/grpc"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

type MemoryService interface {
	GenerateUploadURLs(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, req *dto.GenerateUploadURLsRequest) (*dto.MemoryUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID uuid.UUID, req *dto.ConfirmUploadRequest) error
	GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error)
	ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error)
	DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error
}

type memoryService struct {
	db          *gorm.DB
	memoryRepo  repository.MemoryRepository
	hangoutRepo repository.HangoutRepository
	fileService grpc.FileService
	metrics     *otel.MetricsRecorder
}

func NewMemoryService(db *gorm.DB, memoryRepo repository.MemoryRepository, hangoutRepo repository.HangoutRepository, fileService grpc.FileService, metrics *otel.MetricsRecorder,
) MemoryService {
	return &memoryService{
		db:          db,
		memoryRepo:  memoryRepo,
		hangoutRepo: hangoutRepo,
		fileService: fileService,
		metrics:     metrics,
	}
}

func (s *memoryService) GenerateUploadURLs(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, req *dto.GenerateUploadURLsRequest) (*dto.MemoryUploadResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "memory", "generate_upload")

	ctx, span := otel.StartServiceSpan(ctx, "GenerateUploadURLs",
		attribute.String("user.id", userID.String()),
		attribute.String("hangout.id", hangoutID.String()),
		attribute.Int("file.count", len(req.Files)),
	)
	defer span.End()

	if len(req.Files) > constants.MaxFilePerUpload {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(apperrors.ErrTooManyFiles)
		return nil, apperrors.ErrTooManyFiles
	}

	_, err := s.hangoutRepo.GetHangoutByID(ctx, hangoutID, userID)
	if err != nil {
		recordMetrics("error")
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

		grpcStart := time.Now()
		uploadURLsResp, err = s.fileService.GenerateUploadURLs(ctx, baseStoragePath, fileIntents)
		grpcStatus := "success"
		if err != nil {
			grpcStatus = "error"
		}
		s.metrics.RecordGRPCCall(ctx, "file", "GenerateUploadURLs", grpcStatus, time.Since(grpcStart))

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
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("memory.count", len(memories)))
	span.SetStatusOk()
	recordMetrics("success")
	return mapper.ToMemoryUploadResponse(uploadURLsResp.Urls), nil
}

func (s *memoryService) ConfirmUpload(ctx context.Context, userID uuid.UUID, req *dto.ConfirmUploadRequest) error {
	recordMetrics := s.metrics.StartRequest(ctx, "memory", "confirm_upload")

	ctx, span := otel.StartServiceSpan(ctx, "ConfirmUpload",
		attribute.String("user.id", userID.String()),
		attribute.Int("memory.count", len(req.MemoryIDs)),
	)
	defer span.End()

	memories, err := s.memoryRepo.GetMemoriesByIDs(ctx, req.MemoryIDs, userID)
	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
		return err
	}

	if len(memories) != len(req.MemoryIDs) {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(apperrors.ErrMemoryNotFound)
		return apperrors.ErrMemoryNotFound
	}

	fileIDs := make([]string, 0, len(memories))
	for _, memory := range memories {
		if memory.FileID == nil {
			recordMetrics("error")
			_ = span.RecordErrorWithStatus(apperrors.ErrMemoryNotFound)
			return apperrors.ErrMemoryNotFound
		}
		fileIDs = append(fileIDs, memory.FileID.String())
	}

	grpcStart := time.Now()
	err = s.fileService.ConfirmUpload(ctx, fileIDs)
	grpcStatus := "success"
	if err != nil {
		grpcStatus = "error"
	}
	s.metrics.RecordGRPCCall(ctx, "file", "ConfirmUpload", grpcStatus, time.Since(grpcStart))

	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
	} else {
		span.SetStatusOk()
		recordMetrics("success")
	}
	return err
}

func (s *memoryService) GetMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*dto.MemoryResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "memory", "get")

	ctx, span := otel.StartServiceSpan(ctx, "GetMemory",
		attribute.String("user.id", userID.String()),
		attribute.String("memory.id", memoryID.String()),
	)
	defer span.End()

	memory, err := s.memoryRepo.GetMemoryByID(ctx, memoryID, userID)
	if err != nil {
		recordMetrics("error")
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrMemoryNotFound
		}
		return nil, err
	}

	grpcStart := time.Now()
	fileWithURL, err := s.fileService.GetFileByMemoryID(ctx, memoryID.String())
	grpcStatus := "success"
	if err != nil {
		grpcStatus = "error"
	}
	s.metrics.RecordGRPCCall(ctx, "file", "GetFileByMemoryID", grpcStatus, time.Since(grpcStart))

	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
		return nil, err
	}

	span.SetStatusOk()
	recordMetrics("success")
	return mapper.MemoryToResponseDTO(memory, fileWithURL.DownloadUrl, fileWithURL.FileSize, fileWithURL.MimeType), nil
}

func (s *memoryService) ListMemories(ctx context.Context, userID uuid.UUID, hangoutID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedMemories, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "memory", "list")

	ctx, span := otel.StartServiceSpan(ctx, "ListMemories",
		attribute.String("user.id", userID.String()),
		attribute.String("hangout.id", hangoutID.String()),
		attribute.Int("pagination.limit", pagination.GetLimit()),
	)
	defer span.End()

	_, err := s.hangoutRepo.GetHangoutByID(ctx, hangoutID, userID)
	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrInvalidHangoutID
		}
		return nil, err
	}

	memories, err := s.memoryRepo.GetMemoriesByHangoutID(ctx, hangoutID, pagination)
	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
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

	grpcStart := time.Now()
	filesMap, err := s.fileService.GetFilesByMemoryIDs(ctx, memoryIDs)
	grpcStatus := "success"
	if err != nil {
		grpcStatus = "error"
	}
	s.metrics.RecordGRPCCall(ctx, "file", "GetFilesByMemoryIDs", grpcStatus, time.Since(grpcStart))

	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
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

	span.SetAttributes(
		attribute.Int("memory.count", len(responses)),
		attribute.Bool("pagination.has_more", hasMore),
	)
	span.SetStatusOk()
	recordMetrics("success")
	return &dto.PaginatedMemories{
		Data:       responses,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *memoryService) DeleteMemory(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error {
	recordMetrics := s.metrics.StartRequest(ctx, "memory", "delete")

	ctx, span := otel.StartServiceSpan(ctx, "DeleteMemory",
		attribute.String("user.id", userID.String()),
		attribute.String("memory.id", memoryID.String()),
	)
	defer span.End()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.memoryRepo.WithTx(tx).GetMemoryByID(ctx, memoryID, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return apperrors.ErrMemoryNotFound
			}
			return err
		}

		grpcStart := time.Now()
		deleteErr := s.fileService.DeleteFile(ctx, memoryID.String())
		grpcStatus := "success"
		if deleteErr != nil {
			grpcStatus = "error"
		}
		s.metrics.RecordGRPCCall(ctx, "file", "DeleteFile", grpcStatus, time.Since(grpcStart))

		if deleteErr != nil {
			return deleteErr
		}

		if err := s.memoryRepo.WithTx(tx).DeleteMemory(ctx, memoryID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		recordMetrics("error")
		_ = span.RecordErrorWithStatus(err)
	} else {
		span.SetStatusOk()
		recordMetrics("success")
	}
	return err
}
