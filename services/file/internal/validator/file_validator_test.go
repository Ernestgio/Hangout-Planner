package validator_test

import (
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/validator"
	"github.com/stretchr/testify/require"
)

func TestNewFileValidator(t *testing.T) {
	fv := validator.NewFileValidator()
	require.NotNil(t, fv)
}

func TestFileValidator_ValidateFileUploadIntent(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		size      int64
		mimeType  string
		wantError error
	}{
		{
			name:      "valid jpg file",
			filename:  "photo.jpg",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: nil,
		},
		{
			name:      "valid jpeg file",
			filename:  "photo.jpeg",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: nil,
		},
		{
			name:      "valid png file",
			filename:  "image.png",
			size:      2048,
			mimeType:  "image/png",
			wantError: nil,
		},
		{
			name:      "valid gif file",
			filename:  "animation.gif",
			size:      512,
			mimeType:  "image/gif",
			wantError: nil,
		},
		{
			name:      "valid webp file",
			filename:  "image.webp",
			size:      768,
			mimeType:  "image/webp",
			wantError: nil,
		},
		{
			name:      "uppercase extension",
			filename:  "photo.JPG",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: nil,
		},
		{
			name:      "mixed case extension",
			filename:  "photo.JpEg",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: nil,
		},
		{
			name:      "mime type with charset",
			filename:  "photo.jpg",
			size:      1024,
			mimeType:  "image/jpeg; charset=utf-8",
			wantError: nil,
		},
		{
			name:      "zero size",
			filename:  "photo.jpg",
			size:      0,
			mimeType:  "image/jpeg",
			wantError: apperrors.ErrInvalidFileSize,
		},
		{
			name:      "negative size",
			filename:  "photo.jpg",
			size:      -1,
			mimeType:  "image/jpeg",
			wantError: apperrors.ErrInvalidFileSize,
		},
		{
			name:      "file too large",
			filename:  "photo.jpg",
			size:      constants.MaxFileSize + 1,
			mimeType:  "image/jpeg",
			wantError: apperrors.ErrFileTooLarge,
		},
		{
			name:      "empty filename",
			filename:  "",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: apperrors.ErrInvalidFilename,
		},
		{
			name:      "no extension",
			filename:  "photo",
			size:      1024,
			mimeType:  "image/jpeg",
			wantError: apperrors.ErrInvalidFileExtension,
		},
		{
			name:      "invalid extension",
			filename:  "document.pdf",
			size:      1024,
			mimeType:  "application/pdf",
			wantError: apperrors.ErrInvalidFileExtension,
		},
		{
			name:      "invalid extension txt",
			filename:  "file.txt",
			size:      1024,
			mimeType:  "text/plain",
			wantError: apperrors.ErrInvalidFileExtension,
		},
		{
			name:      "wrong mime type for extension",
			filename:  "photo.jpg",
			size:      1024,
			mimeType:  "image/png",
			wantError: apperrors.ErrInvalidMimeType,
		},
		{
			name:      "invalid mime type",
			filename:  "photo.jpg",
			size:      1024,
			mimeType:  "application/octet-stream",
			wantError: apperrors.ErrInvalidMimeType,
		},
		{
			name:      "empty mime type",
			filename:  "photo.jpg",
			size:      1024,
			mimeType:  "",
			wantError: apperrors.ErrInvalidMimeType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := validator.NewFileValidator()
			err := fv.ValidateFileUploadIntent(tt.filename, tt.size, tt.mimeType)
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFileValidator_GetMaxFileSize(t *testing.T) {
	tests := []struct {
		name     string
		expected int64
	}{
		{
			name:     "returns max file size",
			expected: constants.MaxFileSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := validator.NewFileValidator()
			result := fv.GetMaxFileSize()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFileValidator_IsExtensionAllowed(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		expected  bool
	}{
		{name: "allowed jpg", extension: ".jpg", expected: true},
		{name: "allowed jpeg", extension: ".jpeg", expected: true},
		{name: "allowed png", extension: ".png", expected: true},
		{name: "allowed gif", extension: ".gif", expected: true},
		{name: "allowed webp", extension: ".webp", expected: true},
		{name: "uppercase jpg", extension: ".JPG", expected: true},
		{name: "mixed case jpeg", extension: ".JpEg", expected: true},
		{name: "not allowed pdf", extension: ".pdf", expected: false},
		{name: "not allowed txt", extension: ".txt", expected: false},
		{name: "not allowed doc", extension: ".doc", expected: false},
		{name: "empty extension", extension: "", expected: false},
		{name: "no dot", extension: "jpg", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := validator.NewFileValidator()
			result := fv.IsExtensionAllowed(tt.extension)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFileValidator_IsValidMimeType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		mimeType string
		wantErr  error
	}{
		{name: "valid jpg mime", filename: "photo.jpg", mimeType: "image/jpeg", wantErr: nil},
		{name: "valid jpeg mime", filename: "photo.jpeg", mimeType: "image/jpeg", wantErr: nil},
		{name: "valid png mime", filename: "image.png", mimeType: "image/png", wantErr: nil},
		{name: "valid gif mime", filename: "animation.gif", mimeType: "image/gif", wantErr: nil},
		{name: "valid webp mime", filename: "image.webp", mimeType: "image/webp", wantErr: nil},
		{name: "uppercase mime", filename: "photo.jpg", mimeType: "IMAGE/JPEG", wantErr: nil},
		{name: "mime with whitespace", filename: "photo.jpg", mimeType: "  image/jpeg  ", wantErr: nil},
		{name: "mime with charset", filename: "photo.jpg", mimeType: "image/jpeg; charset=utf-8", wantErr: nil},
		{name: "wrong mime for jpg", filename: "photo.jpg", mimeType: "image/png", wantErr: apperrors.ErrInvalidMimeType},
		{name: "wrong mime for png", filename: "image.png", mimeType: "image/jpeg", wantErr: apperrors.ErrInvalidMimeType},
		{name: "wrong mime for gif", filename: "animation.gif", mimeType: "image/png", wantErr: apperrors.ErrInvalidMimeType},
		{name: "wrong mime for webp", filename: "image.webp", mimeType: "image/jpeg", wantErr: apperrors.ErrInvalidMimeType},
		{name: "invalid mime", filename: "photo.jpg", mimeType: "application/octet-stream", wantErr: apperrors.ErrInvalidMimeType},
		{name: "text mime for image", filename: "photo.png", mimeType: "text/plain", wantErr: apperrors.ErrInvalidMimeType},
		{name: "video mime for image", filename: "photo.jpeg", mimeType: "video/mp4", wantErr: apperrors.ErrInvalidMimeType},
		{name: "empty mime", filename: "photo.jpg", mimeType: "", wantErr: apperrors.ErrInvalidMimeType},
		{name: "mime with spaces gif", filename: "animation.gif", mimeType: "  image/gif  ", wantErr: nil},
		{name: "uppercase mime png", filename: "image.png", mimeType: "IMAGE/PNG", wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := validator.NewFileValidator()
			err := fv.ValidateFileUploadIntent(tt.filename, 1024, tt.mimeType)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
