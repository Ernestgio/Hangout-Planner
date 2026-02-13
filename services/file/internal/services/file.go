package services

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/otel"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/repository"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/storage"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/validator"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
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
	db            *gorm.DB
	fileRepo      repository.MemoryFileRepository
	storage       storage.Storage
	fileValidator validator.FileValidator
	metrics       *otel.MetricsRecorder
}

func NewFileService(db *gorm.DB, repo repository.MemoryFileRepository, storage storage.Storage, fileValidator validator.FileValidator, metrics *otel.MetricsRecorder) FileService {
	return &fileService{
		db:            db,
		fileRepo:      repo,
		storage:       storage,
		fileValidator: fileValidator,
		metrics:       metrics,
	}
}

func (s *fileService) GenerateUploadURLs(ctx context.Context, req *filepb.GenerateUploadURLsRequest) (*filepb.GenerateUploadURLsResponse, error) {
	ctx, span := otel.StartServiceSpan(ctx, "GenerateUploadURLs",
		attribute.Int("file.count", len(req.Files)),
		attribute.String("base_storage_path", req.BaseStoragePath),
	)
	defer span.End()

	recordMetrics := s.metrics.StartOperation(ctx, constants.MetricOpGenerateUploadURL)
	var urls []*filepb.PresignedUploadURL

	err := s.db.Transaction(func(tx *gorm.DB) error {
		files := make([]*domain.MemoryFile, 0, len(req.Files))

		for _, intent := range req.Files {
			if err := s.fileValidator.ValidateFileUploadIntent(intent.Filename, intent.Size, intent.MimeType); err != nil {
				return err
			}

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
			presignedURL := mapper.ToPresignedUploadURL(file.ID, file.MemoryID, file.OriginalName, uploadURL, expiresAt)
			urls = append(urls, presignedURL)
		}

		return nil
	})

	recordMetrics(err)
	if err != nil {
		return nil, span.RecordErrorWithStatus(err)
	}

	for _, intent := range req.Files {
		s.metrics.RecordFileSize(ctx, intent.Size)
	}

	span.SetAttributes(attribute.Int("urls.generated", len(urls)))
	span.SetStatusOk()
	return &filepb.GenerateUploadURLsResponse{
		Urls: urls,
	}, nil
}

func (s *fileService) ConfirmUpload(ctx context.Context, req *filepb.ConfirmUploadRequest) (*filepb.ConfirmUploadResponse, error) {
	ctx, span := otel.StartServiceSpan(ctx, "ConfirmUpload",
		attribute.Int("file.ids.count", len(req.FileIds)),
	)
	defer span.End()

	recordMetrics := s.metrics.StartOperation(ctx, constants.MetricOpConfirmUpload)

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

	recordMetrics(err)
	if err != nil {
		return nil, span.RecordErrorWithStatus(err)
	}

	span.SetStatusOk()
	return &filepb.ConfirmUploadResponse{
		Success: true,
	}, nil
}

func (s *fileService) GetFileByMemoryID(ctx context.Context, req *filepb.GetFileByMemoryIDRequest) (*filepb.GetFileByMemoryIDResponse, error) {
	ctx, span := otel.StartServiceSpan(ctx, "GetFileByMemoryID",
		attribute.String("memory.id", req.MemoryId),
	)
	defer span.End()

	recordMetrics := s.metrics.StartOperation(ctx, constants.MetricOpGetFile)

	memoryID, err := uuid.Parse(req.MemoryId)
	if err != nil {
		recordMetrics(apperrors.ErrInvalidMemoryID)
		return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
	}

	file, err := s.fileRepo.GetByMemoryID(ctx, memoryID)
	if err != nil {
		recordMetrics(apperrors.ErrInvalidMemoryID)
		return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
	}

	downloadURL, err := s.storage.GeneratePresignedDownloadURL(ctx, file.StoragePath)
	if err != nil {
		recordMetrics(err)
		return nil, span.RecordErrorWithStatus(err)
	}

	urlExpiresAt := mapper.GetExpiresAtUnix(s.storage.GetPresignedURLExpiry())
	fileWithURL := mapper.ToFileWithURL(file, downloadURL, urlExpiresAt)

	recordMetrics(nil)
	span.SetAttributes(
		attribute.String("file.id", file.ID.String()),
		attribute.String("file.status", file.FileStatus),
	)
	span.SetStatusOk()
	return &filepb.GetFileByMemoryIDResponse{
		File: fileWithURL,
	}, nil
}

func (s *fileService) GetFilesByMemoryIDs(ctx context.Context, req *filepb.GetFilesByMemoryIDsRequest) (*filepb.GetFilesByMemoryIDsResponse, error) {
	ctx, span := otel.StartServiceSpan(ctx, "GetFilesByMemoryIDs",
		attribute.Int("memory.ids.count", len(req.MemoryIds)),
	)
	defer span.End()

	recordMetrics := s.metrics.StartOperation(ctx, constants.MetricOpGetFilesBatch)

	memoryIDs := make([]uuid.UUID, 0, len(req.MemoryIds))
	for _, idStr := range req.MemoryIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			recordMetrics(apperrors.ErrInvalidMemoryID)
			return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
		}
		memoryIDs = append(memoryIDs, id)
	}

	files, err := s.fileRepo.GetByMemoryIDs(ctx, memoryIDs)
	if err != nil {
		recordMetrics(apperrors.ErrInvalidMemoryID)
		return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
	}

	downloadURLs := make(map[uuid.UUID]string, len(files))
	for _, file := range files {
		downloadURL, err := s.storage.GeneratePresignedDownloadURL(ctx, file.StoragePath)
		if err != nil {
			recordMetrics(err)
			return nil, span.RecordErrorWithStatus(err)
		}
		downloadURLs[file.ID] = downloadURL
	}

	urlExpiresAt := mapper.GetExpiresAtUnix(s.storage.GetPresignedURLExpiry())
	filesWithURLs := mapper.ToFileWithURLBatch(files, downloadURLs, urlExpiresAt)

	recordMetrics(nil)
	span.SetAttributes(attribute.Int("files.returned", len(filesWithURLs)))
	span.SetStatusOk()
	return &filepb.GetFilesByMemoryIDsResponse{
		Files: filesWithURLs,
	}, nil
}

func (s *fileService) DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	ctx, span := otel.StartServiceSpan(ctx, "DeleteFile",
		attribute.String("memory.id", req.MemoryId),
	)
	defer span.End()

	recordMetrics := s.metrics.StartOperation(ctx, constants.MetricOpDeleteFile)

	memoryID, err := uuid.Parse(req.MemoryId)
	if err != nil {
		recordMetrics(apperrors.ErrInvalidMemoryID)
		return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
	}

	file, err := s.fileRepo.GetByMemoryID(ctx, memoryID)
	if err != nil {
		recordMetrics(apperrors.ErrInvalidMemoryID)
		return nil, span.RecordErrorWithStatus(apperrors.ErrInvalidMemoryID)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.fileRepo.WithTx(tx).Delete(ctx, file.ID); err != nil {
			return apperrors.ErrFileDeleteFailed
		}
		return nil
	})

	if err != nil {
		recordMetrics(err)
		return nil, span.RecordErrorWithStatus(err)
	}

	if err := s.storage.Delete(ctx, file.StoragePath); err != nil {
		recordMetrics(apperrors.ErrFileDeleteFailed)
		return nil, span.RecordErrorWithStatus(apperrors.ErrFileDeleteFailed)
	}

	recordMetrics(nil)
	span.SetAttributes(attribute.String("file.id", file.ID.String()))
	span.SetStatusOk()
	return &filepb.DeleteFileResponse{
		Success: true,
	}, nil
}
