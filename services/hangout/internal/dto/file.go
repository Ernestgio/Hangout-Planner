package dto

import (
	"io"

	"github.com/google/uuid"
)

type FileUploadData struct {
	MemoryID    uuid.UUID
	Filename    string
	StoragePath string
	Size        int64
	MimeType    string
	Content     io.Reader
}
