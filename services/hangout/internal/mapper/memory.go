package mapper

import (
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
)

// MemoryToResponseDTO maps domain.Memory + presignedURL to response DTO
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
