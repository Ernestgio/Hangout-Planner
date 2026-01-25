package mapper_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"
	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/mapper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestToFileWithURL(t *testing.T) {
	now := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	memoryID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name         string
		file         *domain.MemoryFile
		downloadURL  string
		urlExpiresAt int64
		expectedID   string
		expectedName string
		expectedSize int64
		expectedMime string
		expectedURL  string
		expectedExp  int64
	}{
		{
			name: "complete file",
			file: &domain.MemoryFile{
				ID:           fileID,
				MemoryID:     memoryID,
				OriginalName: "photo.jpg",
				FileSize:     1024,
				MimeType:     "image/jpeg",
				CreatedAt:    now,
			},
			downloadURL:  "https://s3.example.com/download",
			urlExpiresAt: 1735689599,
			expectedID:   fileID.String(),
			expectedName: "photo.jpg",
			expectedSize: 1024,
			expectedMime: "image/jpeg",
			expectedURL:  "https://s3.example.com/download",
			expectedExp:  1735689599,
		},
		{
			name: "empty values",
			file: &domain.MemoryFile{
				ID:           uuid.Nil,
				MemoryID:     uuid.Nil,
				OriginalName: "",
				FileSize:     0,
				MimeType:     "",
				CreatedAt:    time.Time{},
			},
			downloadURL:  "",
			urlExpiresAt: 0,
			expectedID:   uuid.Nil.String(),
			expectedName: "",
			expectedSize: 0,
			expectedMime: "",
			expectedURL:  "",
			expectedExp:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToFileWithURL(tt.file, tt.downloadURL, tt.urlExpiresAt)
			require.NotNil(t, result)
			require.Equal(t, tt.expectedID, result.Id)
			require.Equal(t, tt.expectedName, result.OriginalName)
			require.Equal(t, tt.expectedSize, result.FileSize)
			require.Equal(t, tt.expectedMime, result.MimeType)
			require.Equal(t, tt.expectedURL, result.DownloadUrl)
			require.Equal(t, tt.expectedExp, result.UrlExpiresAt)
			require.NotNil(t, result.CreatedAt)
		})
	}
}

func TestToFileWithURLBatch(t *testing.T) {
	now := time.Now()
	fileID1 := uuid.New()
	fileID2 := uuid.New()
	memoryID1 := uuid.New()
	memoryID2 := uuid.New()

	tests := []struct {
		name         string
		files        []*domain.MemoryFile
		downloadURLs map[uuid.UUID]string
		urlExpiresAt int64
		expectedLen  int
	}{
		{
			name: "multiple files",
			files: []*domain.MemoryFile{
				{ID: fileID1, MemoryID: memoryID1, OriginalName: "file1.jpg", FileSize: 1024, MimeType: "image/jpeg", CreatedAt: now},
				{ID: fileID2, MemoryID: memoryID2, OriginalName: "file2.png", FileSize: 2048, MimeType: "image/png", CreatedAt: now},
			},
			downloadURLs: map[uuid.UUID]string{
				fileID1: "https://s3.example.com/file1",
				fileID2: "https://s3.example.com/file2",
			},
			urlExpiresAt: 1234567890,
			expectedLen:  2,
		},
		{
			name:         "empty files",
			files:        []*domain.MemoryFile{},
			downloadURLs: map[uuid.UUID]string{},
			urlExpiresAt: 0,
			expectedLen:  0,
		},
		{
			name: "missing download url",
			files: []*domain.MemoryFile{
				{ID: fileID1, MemoryID: memoryID1, OriginalName: "file1.jpg", FileSize: 1024, MimeType: "image/jpeg", CreatedAt: now},
			},
			downloadURLs: map[uuid.UUID]string{},
			urlExpiresAt: 1234567890,
			expectedLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToFileWithURLBatch(tt.files, tt.downloadURLs, tt.urlExpiresAt)
			require.NotNil(t, result)
			require.Len(t, result, tt.expectedLen)
			for _, file := range tt.files {
				fileWithURL, exists := result[file.MemoryID.String()]
				require.True(t, exists)
				require.Equal(t, file.ID.String(), fileWithURL.Id)
				require.Equal(t, file.OriginalName, fileWithURL.OriginalName)
			}
		})
	}
}

func TestToPresignedUploadURL(t *testing.T) {
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	memoryID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name             string
		fileID           uuid.UUID
		memoryID         uuid.UUID
		filename         string
		uploadURL        string
		expiresAt        int64
		expectedFileID   string
		expectedMemoryID string
		expectedFilename string
		expectedURL      string
		expectedExp      int64
	}{
		{
			name:             "complete data",
			fileID:           fileID,
			memoryID:         memoryID,
			filename:         "photo.jpg",
			uploadURL:        "https://s3.example.com/upload",
			expiresAt:        1735689599,
			expectedFileID:   fileID.String(),
			expectedMemoryID: memoryID.String(),
			expectedFilename: "photo.jpg",
			expectedURL:      "https://s3.example.com/upload",
			expectedExp:      1735689599,
		},
		{
			name:             "empty values",
			fileID:           uuid.Nil,
			memoryID:         uuid.Nil,
			filename:         "",
			uploadURL:        "",
			expiresAt:        0,
			expectedFileID:   uuid.Nil.String(),
			expectedMemoryID: uuid.Nil.String(),
			expectedFilename: "",
			expectedURL:      "",
			expectedExp:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToPresignedUploadURL(tt.fileID, tt.memoryID, tt.filename, tt.uploadURL, tt.expiresAt)
			require.NotNil(t, result)
			require.Equal(t, tt.expectedFileID, result.FileId)
			require.Equal(t, tt.expectedMemoryID, result.MemoryId)
			require.Equal(t, tt.expectedFilename, result.Filename)
			require.Equal(t, tt.expectedURL, result.UploadUrl)
			require.Equal(t, tt.expectedExp, result.ExpiresAt)
		})
	}
}

func TestToDomainMemoryFile(t *testing.T) {
	memoryID := uuid.New()
	basePath := "hangouts/123/memories"

	tests := []struct {
		name        string
		intent      *filepb.FileUploadIntent
		basePath    string
		fileStatus  enums.FileUploadStatus
		wantError   bool
		expectedLen int
	}{
		{
			name: "valid intent pending",
			intent: &filepb.FileUploadIntent{
				MemoryId: memoryID.String(),
				Filename: "photo.jpg",
				Size:     1024,
				MimeType: "image/jpeg",
			},
			basePath:   basePath,
			fileStatus: enums.FileUploadStatusPending,
			wantError:  false,
		},
		{
			name: "valid intent completed",
			intent: &filepb.FileUploadIntent{
				MemoryId: memoryID.String(),
				Filename: "document.pdf",
				Size:     2048,
				MimeType: "application/pdf",
			},
			basePath:   basePath,
			fileStatus: enums.FileUploadStatusUploaded,
			wantError:  false,
		},
		{
			name: "invalid memory id",
			intent: &filepb.FileUploadIntent{
				MemoryId: "invalid-uuid",
				Filename: "photo.jpg",
				Size:     1024,
				MimeType: "image/jpeg",
			},
			basePath:   basePath,
			fileStatus: enums.FileUploadStatusPending,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.ToDomainMemoryFile(tt.intent, tt.basePath, tt.fileStatus)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.intent.Filename, result.OriginalName)
				require.Equal(t, tt.intent.Size, result.FileSize)
				require.Equal(t, tt.intent.MimeType, result.MimeType)
				require.Equal(t, string(tt.fileStatus), result.FileStatus)
				require.Contains(t, result.StoragePath, tt.intent.Filename)
			}
		})
	}
}

func TestBuildStoragePath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		memoryID string
		filename string
		expected string
	}{
		{
			name:     "standard path",
			basePath: "hangouts/123/memories",
			memoryID: "456",
			filename: "photo.jpg",
			expected: filepath.Join("hangouts/123/memories", "456", "photo.jpg"),
		},
		{
			name:     "empty base path",
			basePath: "",
			memoryID: "456",
			filename: "photo.jpg",
			expected: filepath.Join("456", "photo.jpg"),
		},
		{
			name:     "all empty",
			basePath: "",
			memoryID: "",
			filename: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.BuildStoragePath(tt.basePath, tt.memoryID, tt.filename)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{name: "jpg extension", filename: "photo.jpg", expected: ".jpg"},
		{name: "png extension", filename: "image.png", expected: ".png"},
		{name: "pdf extension", filename: "document.pdf", expected: ".pdf"},
		{name: "multiple dots", filename: "file.backup.tar.gz", expected: ".gz"},
		{name: "no extension", filename: "filename", expected: ""},
		{name: "empty string", filename: "", expected: ""},
		{name: "dot only", filename: ".", expected: "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GetFileExtension(tt.filename)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetExpiresAtUnix(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{name: "1 hour", duration: 1 * time.Hour},
		{name: "24 hours", duration: 24 * time.Hour},
		{name: "1 minute", duration: 1 * time.Minute},
		{name: "zero duration", duration: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			result := mapper.GetExpiresAtUnix(tt.duration)
			expected := now.Add(tt.duration).Unix()
			require.InDelta(t, expected, result, 2)
		})
	}
}
