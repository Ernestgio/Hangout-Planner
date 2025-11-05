package services

import (
	"context"
	"errors"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, req *dto.CreateActivityRequest) (*dto.ActivityDetailResponse, error)
	GetActivityByID(ctx context.Context, id uuid.UUID) (*dto.ActivityDetailResponse, error)
	GetAllActivities(ctx context.Context) ([]dto.ActivityListItemResponse, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, req *dto.UpdateActivityRequest) (*dto.ActivityDetailResponse, error)
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type activityService struct {
	activityRepo repository.ActivityRepository
	db           *gorm.DB
}

func NewActivityService(db *gorm.DB, activityRepo repository.ActivityRepository) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
		db:           db,
	}
}

func (s *activityService) CreateActivity(ctx context.Context, req *dto.CreateActivityRequest) (*dto.ActivityDetailResponse, error) {
	activityModel := &domain.Activity{
		Name: req.Name,
	}
	createdActivity, err := s.activityRepo.CreateActivity(ctx, activityModel)
	if err != nil {
		return nil, err
	}

	return mapper.ActivitytoDetailResponseDTO(createdActivity, 0), nil
}

func (s *activityService) GetActivityByID(ctx context.Context, id uuid.UUID) (*dto.ActivityDetailResponse, error) {
	activity, hangoutCount, err := s.activityRepo.GetActivityByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ActivitytoDetailResponseDTO(activity, hangoutCount), nil
}

func (s *activityService) GetAllActivities(ctx context.Context) ([]dto.ActivityListItemResponse, error) {
	activitiesWithCount, err := s.activityRepo.GetAllActivities(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ActivityToListItemResponseDTO(activitiesWithCount), nil
}

func (s *activityService) UpdateActivity(ctx context.Context, id uuid.UUID, req *dto.UpdateActivityRequest) (*dto.ActivityDetailResponse, error) {
	var updatedActivity *domain.Activity
	var updatedActivityCount int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.activityRepo.WithTx(tx)

		existingActivity, existingCount, err := txRepo.GetActivityByID(ctx, id)
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
		return nil, err
	}

	return mapper.ActivitytoDetailResponseDTO(updatedActivity, updatedActivityCount), nil
}

func (s *activityService) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.activityRepo.WithTx(tx)
		_, _, err := txRepo.GetActivityByID(ctx, id)
		if err != nil {
			return err
		}
		return txRepo.DeleteActivity(ctx, id)
	})
}
