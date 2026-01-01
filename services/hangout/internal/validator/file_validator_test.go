package validator_test

import (
	"errors"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/validator"
	"github.com/stretchr/testify/require"
)

func TestValidateFileMetadata_TableDriven(t *testing.T) {
	fv := validator.NewFileValidator()

	tests := []struct {
		name     string
		filename string
		size     int64
		mime     string
		wantErr  error
	}{
		{name: "valid jpg", filename: "photo.JPG", size: 1024, mime: "image/jpeg", wantErr: nil},
		{name: "valid png with params", filename: "img.png", size: 1, mime: "image/png; charset=utf-8", wantErr: nil},
		{name: "too large", filename: "big.png", size: constants.MaxFileSize + 1, mime: "image/png", wantErr: apperrors.ErrFileTooLarge},
		{name: "invalid extension", filename: "file.txt", size: 10, mime: "text/plain", wantErr: apperrors.ErrInvalidFileType},
		{name: "mismatched mime", filename: "photo.jpg", size: 10, mime: "image/png", wantErr: apperrors.ErrInvalidFileType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fv.ValidateFileMetadata(tt.filename, tt.size, tt.mime)
			if tt.wantErr == nil {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.True(t, errors.Is(err, tt.wantErr))
		})
	}
}

func TestGetFileExtension_TableDriven(t *testing.T) {
	fv := validator.NewFileValidator()

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{name: "upper ext", filename: "a.JPG", want: ".jpg"},
		{name: "no ext", filename: "file", want: ""},
		{name: "dotfile", filename: ".env", want: ".env"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fv.GetFileExtension(tt.filename)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExtractFileMetadata_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		header   map[string][]string
		filename string
		size     int64
		wantMime string
	}{
		{name: "basic", header: map[string][]string{"Content-Type": {"image/jpeg"}}, filename: "x.jpg", size: 100, wantMime: "image/jpeg"},
		{name: "with params", header: map[string][]string{"Content-Type": {"image/png; charset=utf-8"}}, filename: "y.png", size: 200, wantMime: "image/png; charset=utf-8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &multipart.FileHeader{Filename: tt.filename, Size: tt.size, Header: textproto.MIMEHeader(tt.header)}
			name, size, mime := validator.ExtractFileMetadata(fh)
			require.Equal(t, tt.filename, name)
			require.Equal(t, tt.size, size)
			require.Equal(t, tt.wantMime, mime)
		})
	}
}

func TestIntegration_ExtractAndValidate(t *testing.T) {
	fh := &multipart.FileHeader{Filename: "ok.png", Size: 1024, Header: textproto.MIMEHeader{"Content-Type": {"image/png"}}}
	name, size, mime := validator.ExtractFileMetadata(fh)
	fv := validator.NewFileValidator()
	err := fv.ValidateFileMetadata(name, size, mime)
	require.NoError(t, err)
}
