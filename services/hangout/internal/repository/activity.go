package repository

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	WithTx(tx *gorm.DB) ActivityRepository
	CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	GetActivityByID(ctx context.Context, id uuid.UUID) (*domain.Activity, int64, error)
	GetAllActivities(ctx context.Context) ([]ActivityWithCount, error)
	UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type activityRepository struct {
	db *gorm.DB
}

type ActivityWithCount struct {
	domain.Activity
	HangoutCount int64 `gorm:"column:hangout_count"`
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) WithTx(tx *gorm.DB) ActivityRepository {
	return &activityRepository{db: tx}
}

func (r *activityRepository) CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	if err := r.db.WithContext(ctx).Create(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) GetActivityByID(ctx context.Context, id uuid.UUID) (*domain.Activity, int64, error) {
	var result ActivityWithCount

	err := r.db.WithContext(ctx).
		Model(&domain.Activity{}).
		Select("activities.*, COUNT(hangout_activities.hangout_id) as hangout_count").
		Joins("LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id").
		Where("activities.id = ?", id).
		Group("activities.id").
		First(&result).Error

	if err != nil {
		return nil, 0, err
	}
	activity := result.Activity
	count := result.HangoutCount

	return &activity, count, nil
}

func (r *activityRepository) GetAllActivities(ctx context.Context) ([]ActivityWithCount, error) {
	var results []ActivityWithCount

	err := r.db.WithContext(ctx).
		Model(&domain.Activity{}).
		Select("activities.*, COUNT(hangout_activities.hangout_id) as hangout_count").
		Joins("LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id").
		Group("activities.id").
		Order("activities.name asc").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *activityRepository) UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	if err := r.db.WithContext(ctx).Save(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Activity{}).Error
}
