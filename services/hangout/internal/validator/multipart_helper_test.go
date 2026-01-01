package validator_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/validator"
	"github.com/stretchr/testify/require"
)

func TestMultipartExtractFileMetadata_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		header   textproto.MIMEHeader
		filename string
		size     int64
		wantName string
		wantSize int64
		wantMime string
	}{
		{name: "basic", header: textproto.MIMEHeader{"Content-Type": {"image/jpeg"}}, filename: "a.jpg", size: 123, wantName: "a.jpg", wantSize: 123, wantMime: "image/jpeg"},
		{name: "no header", header: nil, filename: "b.png", size: 0, wantName: "b.png", wantSize: 0, wantMime: ""},
		{name: "multiple values", header: textproto.MIMEHeader{"Content-Type": {"image/png", "charset=utf-8"}}, filename: "c.png", size: 10, wantName: "c.png", wantSize: 10, wantMime: "image/png"},
		{name: "filename path", header: textproto.MIMEHeader{"Content-Type": {"application/octet-stream"}}, filename: "path/to/d.txt", size: 5, wantName: "path/to/d.txt", wantSize: 5, wantMime: "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &multipart.FileHeader{Filename: tt.filename, Size: tt.size, Header: tt.header}
			name, size, mime := validator.ExtractFileMetadata(fh)
			require.Equal(t, tt.wantName, name)
			require.Equal(t, tt.wantSize, size)
			if tt.wantMime == "" {
				require.Equal(t, "", mime)
			} else {
				require.Contains(t, mime, tt.wantMime)
			}
		})
	}
}
