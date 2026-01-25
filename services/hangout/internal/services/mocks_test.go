package services_test

import (
	"context"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, request dto.CreateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, request)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockBcryptUtils struct {
	mock.Mock
}

func (m *MockBcryptUtils) GenerateFromPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockBcryptUtils) CompareHashAndPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

type MockJWTUtils struct {
	mock.Mock
}

func (m *MockJWTUtils) Generate(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) WithTx(tx *gorm.DB) repository.UserRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.UserRepository)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

type MockMemoryRepository struct {
	mock.Mock
}

func (m *MockMemoryRepository) WithTx(tx *gorm.DB) repository.MemoryRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.MemoryRepository)
}

func (m *MockMemoryRepository) CreateMemory(ctx context.Context, memory *domain.Memory) (*domain.Memory, error) {
	args := m.Called(ctx, memory)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Memory), args.Error(1)
}

func (m *MockMemoryRepository) CreateMemoriesBatch(ctx context.Context, memories []*domain.Memory) error {
	args := m.Called(ctx, memories)
	return args.Error(0)
}

func (m *MockMemoryRepository) UpdateFileIDs(ctx context.Context, updates map[uuid.UUID]uuid.UUID) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockMemoryRepository) GetMemoryByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Memory, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Memory), args.Error(1)
}

func (m *MockMemoryRepository) GetMemoriesByIDs(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Memory, error) {
	args := m.Called(ctx, ids, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Memory), args.Error(1)
}

func (m *MockMemoryRepository) GetMemoriesByHangoutID(ctx context.Context, hangoutID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Memory, error) {
	args := m.Called(ctx, hangoutID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Memory), args.Error(1)
}

func (m *MockMemoryRepository) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) GenerateUploadURLs(ctx context.Context, baseStoragePath string, fileIntents []*filepb.FileUploadIntent) (*filepb.GenerateUploadURLsResponse, error) {
	args := m.Called(ctx, baseStoragePath, fileIntents)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*filepb.GenerateUploadURLsResponse), args.Error(1)
}

func (m *MockFileService) ConfirmUpload(ctx context.Context, fileIDs []string) error {
	args := m.Called(ctx, fileIDs)
	return args.Error(0)
}

func (m *MockFileService) GetFileByMemoryID(ctx context.Context, memoryID string) (*filepb.FileWithURL, error) {
	args := m.Called(ctx, memoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*filepb.FileWithURL), args.Error(1)
}

func (m *MockFileService) GetFilesByMemoryIDs(ctx context.Context, memoryIDs []string) (map[string]*filepb.FileWithURL, error) {
	args := m.Called(ctx, memoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*filepb.FileWithURL), args.Error(1)
}

func (m *MockFileService) DeleteFile(ctx context.Context, memoryID string) error {
	args := m.Called(ctx, memoryID)
	return args.Error(0)
}

func (m *MockFileService) Close() error {
	args := m.Called()
	return args.Error(0)
}
