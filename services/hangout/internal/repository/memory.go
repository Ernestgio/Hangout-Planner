package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryRepository interface {
	WithTx(tx *gorm.DB) MemoryRepository
	CreateMemory(ctx context.Context, memory *domain.Memory) (*domain.Memory, error)
	CreateMemoriesBatch(ctx context.Context, memories []*domain.Memory) error
	UpdateFileIDs(ctx context.Context, updates map[uuid.UUID]uuid.UUID) error
	GetMemoryByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Memory, error)
	GetMemoriesByIDs(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Memory, error)
	GetMemoriesByHangoutID(ctx context.Context, hangoutID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Memory, error)
	DeleteMemory(ctx context.Context, id uuid.UUID) error
}

type memoryRepository struct {
	db      *gorm.DB
	metrics *otel.MetricsRecorder
}

func NewMemoryRepository(db *gorm.DB, metrics *otel.MetricsRecorder) MemoryRepository {
	return &memoryRepository{db: db, metrics: metrics}
}

func (r *memoryRepository) WithTx(tx *gorm.DB) MemoryRepository {
	return &memoryRepository{db: tx, metrics: r.metrics}
}

func (r *memoryRepository) CreateMemory(ctx context.Context, memory *domain.Memory) (*domain.Memory, error) {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(memory).Error
	r.metrics.RecordDBOperation(ctx, "insert", "memories", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return memory, nil
}

func (r *memoryRepository) CreateMemoriesBatch(ctx context.Context, memories []*domain.Memory) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(&memories).Error
	r.metrics.RecordDBOperation(ctx, "insert", "memories", time.Since(start), len(memories))
	return err
}

func (r *memoryRepository) UpdateFileIDs(ctx context.Context, updates map[uuid.UUID]uuid.UUID) error {
	if len(updates) == 0 {
		return nil
	}

	ids := make([]interface{}, 0, len(updates))
	caseSQL := "CASE "
	args := make([]interface{}, 0, len(updates)*2)

	for memoryID, fileID := range updates {
		caseSQL += "WHEN id = ? THEN ? "
		args = append(args, memoryID.String(), fileID.String())
		ids = append(ids, memoryID.String())
	}
	caseSQL += "END"

	sql := fmt.Sprintf("UPDATE memories SET file_id = %s WHERE id IN (?%s)", caseSQL, RepeatPlaceholder(len(ids)-1))
	args = append(args, ids...)

	start := time.Now()
	err := r.db.WithContext(ctx).Exec(sql, args...).Error
	r.metrics.RecordDBOperation(ctx, "update", "memories", time.Since(start), len(updates))

	return err
}

func RepeatPlaceholder(count int) string {
	if count <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < count; i++ {
		result += ", ?"
	}
	return result
}

func (r *memoryRepository) GetMemoryByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Memory, error) {
	var memory domain.Memory

	start := time.Now()
	err := r.db.WithContext(ctx).First(&memory, "id = ? AND user_id = ?", id, userID).Error
	r.metrics.RecordDBOperation(ctx, "select", "memories", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return &memory, nil
}

func (r *memoryRepository) GetMemoriesByIDs(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Memory, error) {
	var memories []domain.Memory

	start := time.Now()
	err := r.db.WithContext(ctx).Where("id IN ? AND user_id = ?", ids, userID).Find(&memories).Error
	r.metrics.RecordDBOperation(ctx, "select", "memories", time.Since(start), len(memories))

	if err != nil {
		return nil, err
	}
	return memories, nil
}

func (r *memoryRepository) GetMemoriesByHangoutID(ctx context.Context, hangoutID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Memory, error) {
	start := time.Now()
	var memories []domain.Memory
	limitToFetch := pagination.GetLimit() + 1
	sortByColumn := constants.SortByCreatedAt
	sortDir := pagination.GetSortDir()

	query := r.db.WithContext(ctx).Model(&domain.Memory{}).Where("hangout_id = ?", hangoutID)

	if pagination.AfterID != nil {
		var cursorItem domain.Memory
		if err := r.db.WithContext(ctx).First(&cursorItem, "id = ?", *pagination.AfterID).Error; err != nil {
			return nil, apperrors.ErrInvalidCursorPagination
		}

		cursorValue := cursorItem.CreatedAt

		comparisonOp := ">"
		if sortDir == "desc" {
			comparisonOp = "<"
		}

		query = query.Where(
			fmt.Sprintf("(%s %s ?) OR (%s = ? AND id %s ?)", sortByColumn, comparisonOp, sortByColumn, comparisonOp),
			cursorValue, cursorValue, cursorItem.ID,
		)
	}

	query = query.Order(fmt.Sprintf("%s %s, id %s", sortByColumn, sortDir, sortDir)).Limit(limitToFetch)

	if err := query.Find(&memories).Error; err != nil {
		r.metrics.RecordDBOperation(ctx, "select", "memories", time.Since(start), 0)
		return nil, err
	}

	r.metrics.RecordDBOperation(ctx, "select", "memories", time.Since(start), len(memories))
	return memories, nil
}

func (r *memoryRepository) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Delete(&domain.Memory{}, "id = ?", id).Error
	r.metrics.RecordDBOperation(ctx, "delete", "memories", time.Since(start), 1)
	return err
}
