package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HangoutRepository interface {
	WithTx(tx *gorm.DB) HangoutRepository
	CreateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error)
	GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Hangout, error)
	UpdateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error)
	DeleteHangout(ctx context.Context, id uuid.UUID) error
	GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error)
	GetHangoutActivityIDs(ctx context.Context, hangoutID uuid.UUID) ([]uuid.UUID, error)
	AddHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error
	RemoveHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error
}

type hangoutRepository struct {
	db      *gorm.DB
	metrics *otel.MetricsRecorder
}

func NewHangoutRepository(db *gorm.DB, metrics *otel.MetricsRecorder) HangoutRepository {
	return &hangoutRepository{db: db, metrics: metrics}
}

func (r *hangoutRepository) WithTx(tx *gorm.DB) HangoutRepository {
	return &hangoutRepository{db: tx, metrics: r.metrics}
}

func (r *hangoutRepository) CreateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(hangout).Error
	r.metrics.RecordDBOperation(ctx, "insert", "hangouts", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return hangout, nil
}

func (r *hangoutRepository) GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Hangout, error) {
	var hangout domain.Hangout

	start := time.Now()
	err := r.db.WithContext(ctx).Preload("Activities").First(&hangout, "id = ? AND user_id = ?", id, userID).Error
	r.metrics.RecordDBOperation(ctx, "select", "hangouts", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return &hangout, nil
}

func (r *hangoutRepository) UpdateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	start := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.Hangout{}).
		Where("id = ?", hangout.ID).
		Updates(hangout).Error
	r.metrics.RecordDBOperation(ctx, "update", "hangouts", time.Since(start), 1)

	if err != nil {
		return nil, err
	}

	return hangout, nil
}

func (r *hangoutRepository) DeleteHangout(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.metrics.RecordDBOperation(ctx, "delete", "hangouts", time.Since(start), 0)
		return tx.Error
	}

	if err := tx.Exec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?", id).Error; err != nil {
		tx.Rollback()
		r.metrics.RecordDBOperation(ctx, "delete", "hangouts", time.Since(start), 0)
		return err
	}

	if err := tx.Delete(&domain.Hangout{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		r.metrics.RecordDBOperation(ctx, "delete", "hangouts", time.Since(start), 0)
		return err
	}

	err := tx.Commit().Error
	r.metrics.RecordDBOperation(ctx, "delete", "hangouts", time.Since(start), 1)
	return err
}

func (r *hangoutRepository) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error) {
	start := time.Now()
	var hangouts []domain.Hangout
	limitToFetch := pagination.GetLimit() + 1
	sortByColumn := pagination.GetSortBy()
	sortDir := pagination.GetSortDir()

	query := r.db.WithContext(ctx).Model(&domain.Hangout{}).Where("user_id = ?", userID)

	if pagination.AfterID != nil {
		var cursorItem domain.Hangout
		if err := r.db.WithContext(ctx).First(&cursorItem, "id = ?", *pagination.AfterID).Error; err != nil {
			return nil, apperrors.ErrInvalidCursorPagination
		}

		cursorValue := cursorItem.CreatedAt
		if sortByColumn == constants.SortByDate {
			cursorValue = cursorItem.Date
		}

		comparisonOp := ">"
		if sortDir == constants.SortDirectionDesc {
			comparisonOp = "<"
		}

		query = query.Where(
			fmt.Sprintf("(%s %s ?) OR (%s = ? AND id %s ?)", sortByColumn, comparisonOp, sortByColumn, comparisonOp),
			cursorValue, cursorValue, cursorItem.ID,
		)
	}

	err := query.
		Order(pagination.GetOrderByClause()).
		Limit(limitToFetch).
		Find(&hangouts).Error
	r.metrics.RecordDBOperation(ctx, "select", "hangouts", time.Since(start), len(hangouts))

	if err != nil {
		return nil, err
	}

	return hangouts, nil
}

func (r *hangoutRepository) GetHangoutActivityIDs(ctx context.Context, hangoutID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	start := time.Now()
	err := r.db.WithContext(ctx).
		Table("hangout_activities").
		Where("hangout_id = ?", hangoutID).
		Pluck("activity_id", &ids).Error
	r.metrics.RecordDBOperation(ctx, "select", "hangout_activities", time.Since(start), len(ids))

	return ids, err
}

func (r *hangoutRepository) AddHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error {
	if len(activityIDs) == 0 {
		return nil
	}

	rows := []map[string]interface{}{}
	for _, id := range activityIDs {
		rows = append(rows, map[string]interface{}{
			"hangout_id":  hangoutID,
			"activity_id": id,
		})
	}

	start := time.Now()
	err := r.db.WithContext(ctx).Table("hangout_activities").Create(rows).Error
	r.metrics.RecordDBOperation(ctx, "insert", "hangout_activities", time.Since(start), len(activityIDs))

	return err
}

func (r *hangoutRepository) RemoveHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error {
	if len(activityIDs) == 0 {
		return nil
	}

	start := time.Now()
	err := r.db.WithContext(ctx).
		Table("hangout_activities").
		Where("hangout_id = ? AND activity_id IN ?", hangoutID, activityIDs).
		Delete(nil).Error
	r.metrics.RecordDBOperation(ctx, "delete", "hangout_activities", time.Since(start), len(activityIDs))

	return err
}
