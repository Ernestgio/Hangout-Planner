package dto

import (
	"mime/multipart"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/google/uuid"
)

type CreateMemoriesRequest struct {
	HangoutID uuid.UUID
	Files     []*multipart.FileHeader
}

type MemoryResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	HangoutID uuid.UUID      `json:"hangout_id"`
	FileURL   string         `json:"file_url"`
	FileSize  int64          `json:"file_size"`
	MimeType  string         `json:"mime_type"`
	CreatedAt types.JSONTime `json:"created_at"`
}

type PaginatedMemories struct {
	Data       []MemoryResponse `json:"data"`
	NextCursor *uuid.UUID       `json:"next_cursor"`
	HasMore    bool             `json:"has_more"`
}

type FileUploadIntent struct {
	Filename string `json:"filename" validate:"required"`
	Size     int64  `json:"size" validate:"required,gt=0"`
	MimeType string `json:"mime_type" validate:"required"`
}

type GenerateUploadURLsRequest struct {
	HangoutID uuid.UUID          `json:"hangout_id"`
	Files     []FileUploadIntent `json:"files" validate:"required,dive"`
}

type PresignedUploadURL struct {
	MemoryID  uuid.UUID `json:"memory_id"`
	UploadURL string    `json:"upload_url"`
	ExpiresAt int64     `json:"expires_at"`
}

type MemoryUploadResponse struct {
	UploadURLs []PresignedUploadURL `json:"upload_urls"`
}

type ConfirmUploadRequest struct {
	MemoryIDs []uuid.UUID `json:"memory_ids" validate:"required,dive"`
}
