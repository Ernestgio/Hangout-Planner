package repository

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryFileRepository interface {
	WithTx(tx *gorm.DB) MemoryFileRepository
	CreateBatch(ctx context.Context, files []*domain.MemoryFile) error
	GetByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error)
	GetByMemoryIDs(ctx context.Context, memoryIDs []uuid.UUID) ([]*domain.MemoryFile, error)
	UpdateStatusBatch(ctx context.Context, fileIDs []uuid.UUID, status string) error
	Delete(ctx context.Context, memoryID uuid.UUID) error
}

type memoryFileRepository struct {
	db *gorm.DB
}

func NewMemoryFileRepository(db *gorm.DB) MemoryFileRepository {
	return &memoryFileRepository{db: db}
}

func (r *memoryFileRepository) WithTx(tx *gorm.DB) MemoryFileRepository {
	return &memoryFileRepository{db: tx}
}

func (r *memoryFileRepository) CreateBatch(ctx context.Context, files []*domain.MemoryFile) error {
	return r.db.WithContext(ctx).Create(files).Error
}

func (r *memoryFileRepository) GetByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error) {
	var file domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *memoryFileRepository) GetByMemoryIDs(ctx context.Context, memoryIDs []uuid.UUID) ([]*domain.MemoryFile, error) {
	var files []*domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id IN ?", memoryIDs).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *memoryFileRepository) UpdateStatusBatch(ctx context.Context, fileIDs []uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.MemoryFile{}).
		Where("id IN ?", fileIDs).
		Update("file_status", status).Error
}

func (r *memoryFileRepository) Delete(ctx context.Context, memoryID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("memory_id = ?", memoryID).Delete(&domain.MemoryFile{}).Error
}
