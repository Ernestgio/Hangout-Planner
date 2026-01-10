package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
)

type FileValidator interface {
	ValidateFileUploadIntent(filename string, size int64, mimeType string) error
	GetMaxFileSize() int64
	IsExtensionAllowed(extension string) bool
}

type fileValidator struct {
	maxFileSize       int64
	allowedExtensions map[string]bool
	allowedMimeTypes  map[string][]string
}

func NewFileValidator() FileValidator {
	extensionsMap := make(map[string]bool, len(constants.AllowedFileExtensions))
	for _, ext := range constants.AllowedFileExtensions {
		extensionsMap[strings.ToLower(ext)] = true
	}

	return &fileValidator{
		maxFileSize:       constants.MaxFileSize,
		allowedExtensions: extensionsMap,
		allowedMimeTypes:  constants.AllowedMimeTypes,
	}
}

func (fv *fileValidator) ValidateFileUploadIntent(filename string, size int64, mimeType string) error {
	if size <= 0 {
		return apperrors.ErrInvalidFileSize
	}

	if size > fv.maxFileSize {
		return apperrors.ErrFileTooLarge
	}

	if filename == "" {
		return apperrors.ErrInvalidFilename
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return apperrors.ErrInvalidFileExtension
	}

	if !fv.allowedExtensions[ext] {
		return apperrors.ErrInvalidFileExtension
	}

	if !fv.IsValidMimeType(mimeType, ext) {
		allowedMimes := fv.allowedMimeTypes[ext]
		return fmt.Errorf("%w: MIME type %s doesn't match extension %s. Allowed MIME types: %v",
			apperrors.ErrInvalidMimeType, mimeType, ext, allowedMimes)
	}

	return nil
}

func (fv *fileValidator) IsValidMimeType(mimeType string, ext string) bool {
	allowedMimes, ok := fv.allowedMimeTypes[ext]
	if !ok {
		return false
	}

	mimeTypeLower := strings.ToLower(strings.TrimSpace(mimeType))
	for _, mime := range allowedMimes {
		if strings.HasPrefix(mimeTypeLower, mime) {
			return true
		}
	}

	return false
}

func (fv *fileValidator) GetMaxFileSize() int64 {
	return fv.maxFileSize
}

func (fv *fileValidator) IsExtensionAllowed(extension string) bool {
	return fv.allowedExtensions[strings.ToLower(extension)]
}
