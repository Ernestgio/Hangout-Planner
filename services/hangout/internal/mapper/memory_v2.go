package mapper

import (
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/google/uuid"
)

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
