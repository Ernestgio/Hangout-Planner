package mapper_test

import (
	"testing"
	"time"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/pkg/shared/types"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryToResponseDTO_TableDriven(t *testing.T) {
	now := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name     string
		memory   *domain.Memory
		fileURL  string
		fileSize int64
		mime     string
		wantNil  bool
	}{
		{name: "nil input", memory: nil, wantNil: true},
		{name: "zero time", memory: &domain.Memory{ID: uuid.New(), Name: "m", HangoutID: uuid.New(), CreatedAt: time.Time{}}, fileURL: "", fileSize: 0, mime: "", wantNil: false},
		{name: "with values", memory: &domain.Memory{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Name: "mem1", HangoutID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), CreatedAt: now}, fileURL: "https://x", fileSize: 123, mime: "image/png", wantNil: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.MemoryToResponseDTO(tt.memory, tt.fileURL, tt.fileSize, tt.mime)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, tt.memory.ID, got.ID)
			require.Equal(t, tt.memory.Name, got.Name)
			require.Equal(t, tt.memory.HangoutID, got.HangoutID)
			require.Equal(t, tt.fileURL, got.FileURL)
			require.Equal(t, tt.fileSize, got.FileSize)
			require.Equal(t, tt.mime, got.MimeType)
			require.Equal(t, types.JSONTime(tt.memory.CreatedAt), got.CreatedAt)
		})
	}
}

func TestToMemoryUploadResponse_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		uploadURLs []*filepb.PresignedUploadURL
		wantLen    int
	}{
		{name: "nil input", uploadURLs: nil, wantLen: 0},
		{name: "empty slice", uploadURLs: []*filepb.PresignedUploadURL{}, wantLen: 0},
		{name: "single url", uploadURLs: []*filepb.PresignedUploadURL{{FileId: "f1", MemoryId: "11111111-1111-1111-1111-111111111111", Filename: "a.jpg", UploadUrl: "https://s3/upload1", ExpiresAt: 1234567890}}, wantLen: 1},
		{name: "multiple urls", uploadURLs: []*filepb.PresignedUploadURL{{FileId: "f1", MemoryId: "22222222-2222-2222-2222-222222222222", Filename: "b.png", UploadUrl: "https://s3/upload2", ExpiresAt: 1111111111}, {FileId: "f2", MemoryId: "33333333-3333-3333-3333-333333333333", Filename: "c.pdf", UploadUrl: "https://s3/upload3", ExpiresAt: 2222222222}}, wantLen: 2},
		{name: "invalid uuid", uploadURLs: []*filepb.PresignedUploadURL{{FileId: "f1", MemoryId: "invalid-uuid", Filename: "d.txt", UploadUrl: "https://s3/upload4", ExpiresAt: 3333333333}}, wantLen: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.ToMemoryUploadResponse(tt.uploadURLs)
			require.NotNil(t, got)
			require.Len(t, got.UploadURLs, tt.wantLen)
			for i, url := range tt.uploadURLs {
				require.Equal(t, url.UploadUrl, got.UploadURLs[i].UploadURL)
				require.Equal(t, url.ExpiresAt, got.UploadURLs[i].ExpiresAt)
				expectedID, _ := uuid.Parse(url.MemoryId)
				require.Equal(t, expectedID, got.UploadURLs[i].MemoryID)
			}
		})
	}
}
