package mapper

import (
	"path/filepath"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFileWithURL(file *domain.MemoryFile, downloadURL string, urlExpiresAt int64) *filepb.FileWithURL {
	return &filepb.FileWithURL{
		Id:           file.ID.String(),
		OriginalName: file.OriginalName,
		FileSize:     file.FileSize,
		MimeType:     file.MimeType,
		CreatedAt:    timestamppb.New(file.CreatedAt),
		DownloadUrl:  downloadURL,
		UrlExpiresAt: urlExpiresAt,
	}
}

func ToFileWithURLBatch(files []*domain.MemoryFile, downloadURLs map[uuid.UUID]string, urlExpiresAt int64) map[string]*filepb.FileWithURL {
	result := make(map[string]*filepb.FileWithURL, len(files))
	for _, file := range files {
		downloadURL := downloadURLs[file.ID]
		result[file.MemoryID.String()] = ToFileWithURL(file, downloadURL, urlExpiresAt)
	}
	return result
}

func ToPresignedUploadURL(fileID uuid.UUID, memoryID uuid.UUID, filename, uploadURL string, expiresAt int64) *filepb.PresignedUploadURL {
	return &filepb.PresignedUploadURL{
		FileId:    fileID.String(),
		MemoryId:  memoryID.String(),
		Filename:  filename,
		UploadUrl: uploadURL,
		ExpiresAt: expiresAt,
	}
}

func ToDomainMemoryFile(intent *filepb.FileUploadIntent, basePath string, fileStatus enums.FileUploadStatus) (*domain.MemoryFile, error) {
	memoryID, err := uuid.Parse(intent.MemoryId)
	if err != nil {
		return nil, err
	}

	storagePath := BuildStoragePath(basePath, intent.MemoryId, intent.Filename)

	return &domain.MemoryFile{
		MemoryID:     memoryID,
		OriginalName: intent.Filename,
		StoragePath:  storagePath,
		FileSize:     intent.Size,
		MimeType:     intent.MimeType,
		FileStatus:   string(fileStatus),
	}, nil
}

func BuildStoragePath(basePath, memoryID, filename string) string {
	return filepath.Join(basePath, memoryID, filename)
}

func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

func GetExpiresAtUnix(duration time.Duration) int64 {
	return time.Now().Add(duration).Unix()
}
