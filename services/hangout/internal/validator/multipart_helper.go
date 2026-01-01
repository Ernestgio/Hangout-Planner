package validator

import (
	"mime/multipart"
)

func ExtractFileMetadata(file *multipart.FileHeader) (filename string, size int64, mimeType string) {
	return file.Filename, file.Size, file.Header.Get("Content-Type")
}
