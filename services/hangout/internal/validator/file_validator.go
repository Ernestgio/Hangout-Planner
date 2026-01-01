package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
)

type FileValidator struct {
	maxFileSize       int64
	allowedExtensions map[string]bool
}

func NewFileValidator() *FileValidator {
	extensions := strings.Split(constants.AllowedImageExtension, ",")
	allowedMap := make(map[string]bool)
	for _, ext := range extensions {
		allowedMap[strings.ToLower(strings.TrimSpace(ext))] = true
	}

	return &FileValidator{
		maxFileSize:       constants.MaxFileSize,
		allowedExtensions: allowedMap,
	}
}

func (fv *FileValidator) ValidateFileMetadata(filename string, size int64, mimeType string) error {
	if size > fv.maxFileSize {
		return fmt.Errorf("%w: file size %d exceeds maximum %d bytes", apperrors.ErrFileTooLarge, size, fv.maxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !fv.allowedExtensions[ext] {
		return fmt.Errorf("%w: extension %s not allowed. Allowed: %s", apperrors.ErrInvalidFileType, ext, constants.AllowedImageExtension)
	}

	if !fv.isValidMimeType(mimeType, ext) {
		return fmt.Errorf("%w: MIME type %s doesn't match extension %s", apperrors.ErrInvalidFileType, mimeType, ext)
	}

	return nil
}

func (fv *FileValidator) isValidMimeType(mimeType string, ext string) bool {
	validMimes := map[string][]string{
		".jpg":  {"image/jpeg", "image/jpg"},
		".jpeg": {"image/jpeg", "image/jpg"},
		".png":  {"image/png"},
		".gif":  {"image/gif"},
		".webp": {"image/webp"},
	}

	allowedMimes, ok := validMimes[ext]
	if !ok {
		return false
	}

	for _, mime := range allowedMimes {
		if strings.HasPrefix(strings.ToLower(mimeType), mime) {
			return true
		}
	}

	return false
}

func (fv *FileValidator) GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}
