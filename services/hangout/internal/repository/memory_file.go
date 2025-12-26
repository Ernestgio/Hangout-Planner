package repository

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryFileRepository interface {
	WithTx(tx *gorm.DB) MemoryFileRepository
	CreateFile(ctx context.Context, file *domain.MemoryFile) (*domain.MemoryFile, error)
	GetFileByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error)
	DeleteFile(ctx context.Context, memoryID uuid.UUID) error
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

func (r *memoryFileRepository) CreateFile(ctx context.Context, file *domain.MemoryFile) (*domain.MemoryFile, error) {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *memoryFileRepository) GetFileByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error) {
	var file domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *memoryFileRepository) DeleteFile(ctx context.Context, memoryID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("memory_id = ?", memoryID).Delete(&domain.MemoryFile{}).Error
}
