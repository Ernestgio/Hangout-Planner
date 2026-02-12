package services

import (
	"context"
	"errors"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, userID uuid.UUID, req *dto.CreateActivityRequest) (*dto.ActivityDetailResponse, error)
	GetActivityByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ActivityDetailResponse, error)
	GetAllActivities(ctx context.Context, userID uuid.UUID) ([]dto.ActivityListItemResponse, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateActivityRequest) (*dto.ActivityDetailResponse, error)
	DeleteActivity(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type activityService struct {
	activityRepo repository.ActivityRepository
	db           *gorm.DB
	metrics      *otel.MetricsRecorder
}

func NewActivityService(db *gorm.DB, activityRepo repository.ActivityRepository, metrics *otel.MetricsRecorder) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
		db:           db,
		metrics:      metrics,
	}
}

func (s *activityService) CreateActivity(ctx context.Context, userID uuid.UUID, req *dto.CreateActivityRequest) (*dto.ActivityDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "activity", "create")

	activityModel := &domain.Activity{
		Name: req.Name,
	}
	activityModel.UserID = &userID

	createdActivity, err := s.activityRepo.CreateActivity(ctx, activityModel)
	if err != nil {
		recordMetrics("error")
		return nil, err
	}

	recordMetrics("success")
	return mapper.ActivitytoDetailResponseDTO(createdActivity, 0), nil
}

func (s *activityService) GetActivityByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ActivityDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "activity", "get")

	activity, hangoutCount, err := s.activityRepo.GetActivityByID(ctx, id, userID)
	if err != nil {
		recordMetrics("error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	recordMetrics("success")
	return mapper.ActivitytoDetailResponseDTO(activity, hangoutCount), nil
}

func (s *activityService) GetAllActivities(ctx context.Context, userID uuid.UUID) ([]dto.ActivityListItemResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "activity", "list")

	activitiesWithCount, err := s.activityRepo.GetAllActivities(ctx, userID)
	if err != nil {
		recordMetrics("error")
		return nil, err
	}

	recordMetrics("success")
	return mapper.ActivityToListItemResponseDTO(activitiesWithCount), nil
}

func (s *activityService) UpdateActivity(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateActivityRequest) (*dto.ActivityDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "activity", "update")

	var updatedActivity *domain.Activity
	var updatedActivityCount int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.activityRepo.WithTx(tx)

		existingActivity, existingCount, err := txRepo.GetActivityByID(ctx, id, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return err
		}

		mapper.ApplyUpdateToActivity(existingActivity, req)
		updatedActivityCount = existingCount
		updatedActivity, err = txRepo.UpdateActivity(ctx, existingActivity)

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		recordMetrics("error")
		return nil, err
	}

	recordMetrics("success")
	return mapper.ActivitytoDetailResponseDTO(updatedActivity, updatedActivityCount), nil
}

func (s *activityService) DeleteActivity(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	recordMetrics := s.metrics.StartRequest(ctx, "activity", "delete")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.activityRepo.WithTx(tx)
		_, _, err := txRepo.GetActivityByID(ctx, id, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return err
		}
		return txRepo.DeleteActivity(ctx, id)
	})

	if err != nil {
		recordMetrics("error")
	} else {
		recordMetrics("success")
	}
	return err
}
