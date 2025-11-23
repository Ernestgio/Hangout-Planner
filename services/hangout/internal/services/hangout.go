package services

import (
	"context"
	"errors"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/enums"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HangoutService interface {
	CreateHangout(ctx context.Context, userID uuid.UUID, req *dto.CreateHangoutRequest) (*dto.HangoutDetailResponse, error)
	GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.HangoutDetailResponse, error)
	UpdateHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateHangoutRequest) (*dto.HangoutDetailResponse, error)
	DeleteHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedHangouts, error)
}

type hangoutService struct {
	db           *gorm.DB
	hangoutRepo  repository.HangoutRepository
	activityRepo repository.ActivityRepository
}

func NewHangoutService(db *gorm.DB, hangoutRepo repository.HangoutRepository, activityRepo repository.ActivityRepository) HangoutService {
	return &hangoutService{
		db:           db,
		hangoutRepo:  hangoutRepo,
		activityRepo: activityRepo,
	}
}

func (s *hangoutService) CreateHangout(ctx context.Context, userID uuid.UUID, req *dto.CreateHangoutRequest) (*dto.HangoutDetailResponse, error) {
	hangoutModel, err := mapper.HangoutCreateRequestToModel(req)
	if err != nil {
		return nil, err
	}

	if req.Status == "" {
		hangoutModel.Status = enums.StatusPlanning
	}

	if len(req.ActivityIDs) > 0 {
		activities, err := s.activityRepo.GetActivitiesByIDs(ctx, req.ActivityIDs)
		if err != nil {
			return nil, err
		}

		hangoutModel.Activities = activities
	}

	hangoutModel.UserID = &userID
	createdHangout, err := s.hangoutRepo.CreateHangout(ctx, hangoutModel)
	if err != nil {
		return nil, err
	}

	return mapper.HangoutToDetailResponseDTO(createdHangout), nil
}

func (s *hangoutService) GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.HangoutDetailResponse, error) {
	hangout, err := s.hangoutRepo.GetHangoutByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return mapper.HangoutToDetailResponseDTO(hangout), nil
}

func (s *hangoutService) UpdateHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateHangoutRequest) (*dto.HangoutDetailResponse, error) {
	var updatedHangout *domain.Hangout

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txHangoutRepo := s.hangoutRepo.WithTx(tx)
		txActivityRepo := s.activityRepo.WithTx(tx)

		existingHangout, err := txHangoutRepo.GetHangoutByID(ctx, id, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return err
		}

		err = mapper.ApplyUpdateToHangout(existingHangout, req)
		if err != nil {
			return err
		}

		ids := req.ActivityIDs
		if len(ids) > 0 {
			existingActivities, err := txActivityRepo.GetActivitiesByIDs(ctx, ids)
			if err != nil {
				return err
			}
			if len(existingActivities) != len(ids) {
				return apperrors.ErrInvalidActivityIDs
			}
		}

		if err := txHangoutRepo.ReplaceHangoutActivities(ctx, existingHangout.ID, ids); err != nil {
			return err
		}

		_, err = txHangoutRepo.UpdateHangout(ctx, existingHangout)
		if err != nil {
			return err
		}

		updatedHangout, err = txHangoutRepo.GetHangoutByID(ctx, id, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mapper.HangoutToDetailResponseDTO(updatedHangout), nil
}

func (s *hangoutService) DeleteHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.hangoutRepo.WithTx(tx)
		_, err := txRepo.GetHangoutByID(ctx, id, userID)
		if err != nil {
			return err
		}
		return txRepo.DeleteHangout(ctx, id)
	})
}

func (s *hangoutService) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedHangouts, error) {
	hangouts, err := s.hangoutRepo.GetHangoutsByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	var nextCursor *uuid.UUID
	hasMore := false
	limit := pagination.GetLimit()

	if len(hangouts) > limit {
		hasMore = true
		nextCursor = &hangouts[limit-1].ID
		hangouts = hangouts[:limit]
	}

	responseDTOs := mapper.HangoutsToListItemResponseDTOs(hangouts)

	return &dto.PaginatedHangouts{
		Data:       responseDTOs,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil

}
