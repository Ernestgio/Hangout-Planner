package repository

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/otel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := otel.StartRepositorySpan(ctx, "CreateBatch",
		attribute.String("db.operation", "insert"),
		attribute.String("db.table", "memory_files"),
		attribute.Int("db.batch_size", len(files)),
	)
	defer span.End()

	start := time.Now()
	err := r.db.WithContext(ctx).Create(files).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpInsert, time.Since(start), len(files))

	if err != nil {
		return span.RecordErrorWithStatus(err)
	}

	span.SetStatusOk()
	return nil
}

func (r *memoryFileRepository) GetByMemoryID(ctx context.Context, memoryID uuid.UUID) (*domain.MemoryFile, error) {
	ctx, span := otel.StartRepositorySpan(ctx, "GetByMemoryID",
		attribute.String("db.operation", "select"),
		attribute.String("db.table", "memory_files"),
		attribute.String("memory.id", memoryID.String()),
	)
	defer span.End()

	start := time.Now()
	var file domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&file).Error; err != nil {
		r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), 1)
		return nil, span.RecordErrorWithStatus(err)
	}
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), 1)
	span.SetStatusOk()
	return &file, nil
}

func (r *memoryFileRepository) GetByMemoryIDs(ctx context.Context, memoryIDs []uuid.UUID) ([]*domain.MemoryFile, error) {
	ctx, span := otel.StartRepositorySpan(ctx, "GetByMemoryIDs",
		attribute.String("db.operation", "select"),
		attribute.String("db.table", "memory_files"),
		attribute.Int("memory.ids.count", len(memoryIDs)),
	)
	defer span.End()

	start := time.Now()
	var files []*domain.MemoryFile
	if err := r.db.WithContext(ctx).Where("memory_id IN ?", memoryIDs).Find(&files).Error; err != nil {
		r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), len(memoryIDs))
		return nil, span.RecordErrorWithStatus(err)
	}
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpSelect, time.Since(start), len(memoryIDs))
	span.SetAttributes(attribute.Int("files.found", len(files)))
	span.SetStatusOk()
	return files, nil
}

func (r *memoryFileRepository) UpdateStatusBatch(ctx context.Context, fileIDs []uuid.UUID, status string) error {
	ctx, span := otel.StartRepositorySpan(ctx, "UpdateStatusBatch",
		attribute.String("db.operation", "update"),
		attribute.String("db.table", "memory_files"),
		attribute.Int("file.ids.count", len(fileIDs)),
		attribute.String("file.status", status),
	)
	defer span.End()

	start := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.MemoryFile{}).
		Where("id IN ?", fileIDs).
		Update("file_status", status).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpUpdate, time.Since(start), len(fileIDs))

	if err != nil {
		return span.RecordErrorWithStatus(err)
	}

	span.SetStatusOk()
	return nil
}

func (r *memoryFileRepository) Delete(ctx context.Context, memoryID uuid.UUID) error {
	ctx, span := otel.StartRepositorySpan(ctx, "Delete",
		attribute.String("db.operation", "delete"),
		attribute.String("db.table", "memory_files"),
		attribute.String("memory.id", memoryID.String()),
	)
	defer span.End()

	start := time.Now()
	err := r.db.WithContext(ctx).Where("memory_id = ?", memoryID).Delete(&domain.MemoryFile{}).Error
	r.metrics.RecordDBOperation(ctx, constants.MetricDBOpDelete, time.Since(start), 1)

	if err != nil {
		return span.RecordErrorWithStatus(err)
	}

	span.SetStatusOk()
	return nil
}
