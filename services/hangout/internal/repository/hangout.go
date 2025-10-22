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
	if err := r.db.WithContext(ctx).First(&hangout, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &hangout, nil
}

func (r *hangoutRepository) UpdateHangout(ctx context.Context, hangout *domain.Hangout) (*domain.Hangout, error) {
	if err := r.db.WithContext(ctx).Save(hangout).Error; err != nil {
		return nil, err
	}
	return hangout, nil
}
func (r *hangoutRepository) DeleteHangout(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Hangout{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
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
