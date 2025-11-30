package repository

import (
	"context"
	"fmt"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
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
	db *gorm.DB
}

func NewHangoutRepository(db *gorm.DB) HangoutRepository {
	return &hangoutRepository{db: db}
}

func (r *hangoutRepository) WithTx(tx *gorm.DB) HangoutRepository {
	return &hangoutRepository{db: tx}
}

func (r *hangoutRepository) CreateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	if err := r.db.WithContext(ctx).Create(hangout).Error; err != nil {
		return nil, err
	}
	return hangout, nil
}

func (r *hangoutRepository) GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Hangout, error) {
	var hangout domain.Hangout
	if err := r.db.WithContext(ctx).Preload("Activities").First(&hangout, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &hangout, nil
}

func (r *hangoutRepository) UpdateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	if err := r.db.WithContext(ctx).
		Model(&domain.Hangout{}).
		Where("id = ?", hangout.ID).
		Updates(hangout).Error; err != nil {
		return nil, err
	}

	return hangout, nil
}

func (r *hangoutRepository) DeleteHangout(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Exec("DELETE FROM `hangout_activities` WHERE `hangout_id` = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&domain.Hangout{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *hangoutRepository) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error) {
	var hangouts []domain.Hangout
	limitToFetch := pagination.GetLimit() + 1
	sortByColumn := pagination.GetSortBy()
	sortDir := pagination.GetSortDir()

	query := r.db.WithContext(ctx).Model(&domain.Hangout{}).Where("user_id = ?", userID)

	if pagination.AfterID != nil {
		var cursorItem domain.Hangout
		if err := r.db.WithContext(ctx).First(&cursorItem, "id = ?", *pagination.AfterID).Error; err != nil {
			return nil, fmt.Errorf("cursor item not found: %w", err)
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

	if err != nil {
		return nil, err
	}

	return hangouts, nil
}

func (r *hangoutRepository) GetHangoutActivityIDs(ctx context.Context, hangoutID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Table("hangout_activities").
		Where("hangout_id = ?", hangoutID).
		Pluck("activity_id", &ids).Error

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
	return r.db.WithContext(ctx).Table("hangout_activities").Create(rows).Error
}

func (r *hangoutRepository) RemoveHangoutActivities(ctx context.Context, hangoutID uuid.UUID, activityIDs []uuid.UUID) error {
	if len(activityIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Table("hangout_activities").
		Where("hangout_id = ? AND activity_id IN ?", hangoutID, activityIDs).
		Delete(nil).Error
}
