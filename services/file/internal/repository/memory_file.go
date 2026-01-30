package repository

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/otel"
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
	db      *gorm.DB
	metrics *otel.MetricsRecorder
}

func NewMemoryFileRepository(db *gorm.DB, metrics *otel.MetricsRecorder) MemoryFileRepository {
	return &memoryFileRepository{db: db, metrics: metrics}
}

func (r *memoryFileRepository) WithTx(tx *gorm.DB) MemoryFileRepository {
	return &memoryFileRepository{db: tx, metrics: r.metrics}
}

func (r *memoryFileRepository) CreateBatch(ctx context.Context, files []*domain.MemoryFile) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(files).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpInsert, time.Since(start), len(files))
	return err
}

func (r *memoryFileRepository) GetByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error) {
	start := time.Now()
	var file domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&file).Error; err != nil {
		r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), 1)
		return nil, err
	}
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), 1)
	return &file, nil
}

func (r *memoryFileRepository) GetByMemoryIDs(ctx context.Context, memoryIDs []uuid.UUID) ([]*domain.MemoryFile, error) {
	start := time.Now()
	var files []*domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id IN ?", memoryIDs).Find(&files).Error; err != nil {
		r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), len(memoryIDs))
		return nil, err
	}
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), len(memoryIDs))
	return files, nil
}

func (r *memoryFileRepository) UpdateStatusBatch(ctx context.Context, fileIDs []uuid.UUID, status string) error {
	start := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.MemoryFile{}).
		Where("id IN ?", fileIDs).
		Update("file_status", status).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpUpdate, time.Since(start), len(fileIDs))
	return err
}

func (r *memoryFileRepository) Delete(ctx context.Context, memoryID uuid.UUID) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).Delete(&domain.MemoryFile{}).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpDelete, time.Since(start), 1)
	return err
}
