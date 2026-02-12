package repository

import (
	"context"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	WithTx(tx *gorm.DB) ActivityRepository
	CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	GetActivityByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Activity, int64, error)
	GetAllActivities(ctx context.Context, userID uuid.UUID) ([]ActivityWithCount, error)
	UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	DeleteActivity(ctx context.Context, id uuid.UUID) error
	GetActivitiesByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Activity, error)
}

type activityRepository struct {
	db      *gorm.DB
	metrics *otel.MetricsRecorder
}

type ActivityWithCount struct {
	domain.Activity
	HangoutCount int64 `gorm:"column:hangout_count"`
}

func NewActivityRepository(db *gorm.DB, metrics *otel.MetricsRecorder) ActivityRepository {
	return &activityRepository{db: db, metrics: metrics}
}

func (r *activityRepository) WithTx(tx *gorm.DB) ActivityRepository {
	return &activityRepository{db: tx, metrics: r.metrics}
}

func (r *activityRepository) CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(activity).Error
	r.metrics.RecordDBOperation(ctx, "insert", "activities", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) GetActivityByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Activity, int64, error) {
	var result ActivityWithCount

	start := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.Activity{}).
		Select("activities.*, COUNT(hangout_activities.hangout_id) as hangout_count").
		Joins("LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id").
		Where("activities.id = ?", id).
		Where("activities.user_id = ?", userID).
		Group("activities.id").
		First(&result).Error
	r.metrics.RecordDBOperation(ctx, "select", "activities", time.Since(start), 1)

	if err != nil {
		return nil, 0, err
	}
	activity := result.Activity
	count := result.HangoutCount

	return &activity, count, nil
}

func (r *activityRepository) GetAllActivities(ctx context.Context, userID uuid.UUID) ([]ActivityWithCount, error) {
	var results []ActivityWithCount

	start := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.Activity{}).
		Select("activities.*, COUNT(hangout_activities.hangout_id) as hangout_count").
		Where("activities.user_id = ?", userID).
		Joins("LEFT JOIN hangout_activities ON hangout_activities.activity_id = activities.id").
		Group("activities.id").
		Order("activities.name asc").
		Find(&results).Error
	r.metrics.RecordDBOperation(ctx, "select", "activities", time.Since(start), len(results))

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *activityRepository) UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	start := time.Now()
	err := r.db.WithContext(ctx).Save(activity).Error
	r.metrics.RecordDBOperation(ctx, "update", "activities", time.Since(start), 1)

	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Activity{}).Error
	r.metrics.RecordDBOperation(ctx, "delete", "activities", time.Since(start), 1)
	return err
}

func (r *activityRepository) GetActivitiesByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Activity, error) {
	var activities []*domain.Activity

	start := time.Now()
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&activities).Error
	r.metrics.RecordDBOperation(ctx, "select", "activities", time.Since(start), len(activities))

	if err != nil {
		return nil, err
	}
	return activities, nil
}
