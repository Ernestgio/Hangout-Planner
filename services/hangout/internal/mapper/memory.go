package mapper

import (
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/google/uuid"
)

func MemoryToResponseDTO(memory *domain.Memory, fileURL string, fileSize int64, mimeType string) *dto.MemoryResponse {
	if memory == nil {
		return nil
	}

	return &dto.MemoryResponse{
		ID:        memory.ID,
		Name:      memory.Name,
		HangoutID: memory.HangoutID,
		FileURL:   fileURL,
		FileSize:  fileSize,
		MimeType:  mimeType,
		CreatedAt: types.JSONTime(memory.CreatedAt),
	}
}

func ToMemoryUploadResponse(uploadURLs []*filepb.PresignedUploadURL) *dto.MemoryUploadResponse {
	urls := make([]dto.PresignedUploadURL, len(uploadURLs))

	for i, url := range uploadURLs {
		memoryID, _ := uuid.Parse(url.MemoryId)
		urls[i] = dto.PresignedUploadURL{
			MemoryID:  memoryID,
			UploadURL: url.UploadUrl,
			ExpiresAt: url.ExpiresAt,
		}
	}

	return &dto.MemoryUploadResponse{
		UploadURLs: urls,
	}
}
