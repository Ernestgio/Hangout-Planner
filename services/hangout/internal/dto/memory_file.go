package dto

import "github.com/google/uuid"

type MemoryFileResponse struct {
	ID            uuid.UUID `json:"id"`
	MemoryID      uuid.UUID `json:"memory_id"`
	OriginalName  string    `json:"original_name"`
	FileExtension string    `json:"file_extension"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	FileURL       string    `json:"file_url"`
}
