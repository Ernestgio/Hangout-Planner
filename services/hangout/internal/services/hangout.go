package services

import (
	"context"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
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
	metrics      *otel.MetricsRecorder
}

func NewHangoutService(db *gorm.DB, hangoutRepo repository.HangoutRepository, activityRepo repository.ActivityRepository, metrics *otel.MetricsRecorder) HangoutService {
	return &hangoutService{
		db:           db,
		hangoutRepo:  hangoutRepo,
		activityRepo: activityRepo,
		metrics:      metrics,
	}
}

func (s *hangoutService) CreateHangout(ctx context.Context, userID uuid.UUID, req *dto.CreateHangoutRequest) (*dto.HangoutDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "hangout", "create")

	hangoutModel, err := mapper.HangoutCreateRequestToModel(req)
	if err != nil {
		recordMetrics("error")
		return nil, err
	}
	hangoutModel.UserID = &userID

	var created *domain.Hangout
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txHangoutRepo := s.hangoutRepo.WithTx(tx)
		txActivityRepo := s.activityRepo.WithTx(tx)

		if len(req.ActivityIDs) > 0 {
			acts, err := txActivityRepo.GetActivitiesByIDs(ctx, req.ActivityIDs)
			if err != nil {
				return err
			}

			if len(acts) != len(req.ActivityIDs) {
				return apperrors.ErrInvalidActivityIDs
			}
		}

		h, err := txHangoutRepo.CreateHangout(ctx, hangoutModel)
		if err != nil {
			return err
		}
		created = h

		if len(req.ActivityIDs) > 0 {
			if err := txHangoutRepo.AddHangoutActivities(ctx, created.ID, req.ActivityIDs); err != nil {
				return err
			}
		}

		created, err = txHangoutRepo.GetHangoutByID(ctx, created.ID, userID)
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
	return mapper.HangoutToDetailResponseDTO(created), nil
}

func (s *hangoutService) GetHangoutByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.HangoutDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "hangout", "get")

	hangout, err := s.hangoutRepo.GetHangoutByID(ctx, id, userID)
	if err != nil {
		recordMetrics("error")
		return nil, err
	}

	recordMetrics("success")
	return mapper.HangoutToDetailResponseDTO(hangout), nil
}

func (s *hangoutService) UpdateHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateHangoutRequest) (*dto.HangoutDetailResponse, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "hangout", "update")

	var updatedHangout *domain.Hangout

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txHangoutRepo := s.hangoutRepo.WithTx(tx)
		txActivityRepo := s.activityRepo.WithTx(tx)

		existingHangout, err := txHangoutRepo.GetHangoutByID(ctx, id, userID)
		if err != nil {
			return err
		}

		err = mapper.ApplyUpdateToHangout(existingHangout, req)
		if err != nil {
			return err
		}

		if req.ActivityIDs != nil {
			acts, err := txActivityRepo.GetActivitiesByIDs(ctx, req.ActivityIDs)
			if err != nil {
				return err
			}

			if len(acts) != len(req.ActivityIDs) {
				return apperrors.ErrInvalidActivityIDs
			}
		}

		_, err = txHangoutRepo.UpdateHangout(ctx, existingHangout)
		if err != nil {
			return err
		}

		if req.ActivityIDs != nil {
			newIDs := req.ActivityIDs

			currentMap := make(map[uuid.UUID]bool)
			for _, act := range existingHangout.Activities {
				currentMap[act.ID] = true
			}

			newMap := make(map[uuid.UUID]bool)
			for _, id := range newIDs {
				newMap[id] = true
			}

			var toRemove []uuid.UUID
			for id := range currentMap {
				if !newMap[id] {
					toRemove = append(toRemove, id)
				}
			}

			var toAdd []uuid.UUID
			for id := range newMap {
				if !currentMap[id] {
					toAdd = append(toAdd, id)
				}
			}

			if len(toRemove) > 0 {
				if err := txHangoutRepo.RemoveHangoutActivities(ctx, existingHangout.ID, toRemove); err != nil {
					return err
				}
			}

			if len(toAdd) > 0 {
				if err := txHangoutRepo.AddHangoutActivities(ctx, existingHangout.ID, toAdd); err != nil {
					return err
				}
			}
		}

		updatedHangout, err = txHangoutRepo.GetHangoutByID(ctx, id, userID)
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
	return mapper.HangoutToDetailResponseDTO(updatedHangout), nil
}

func (s *hangoutService) DeleteHangout(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	recordMetrics := s.metrics.StartRequest(ctx, "hangout", "delete")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.hangoutRepo.WithTx(tx)
		_, err := txRepo.GetHangoutByID(ctx, id, userID)
		if err != nil {
			return err
		}
		return txRepo.DeleteHangout(ctx, id)
	})

	if err != nil {
		recordMetrics("error")
	} else {
		recordMetrics("success")
	}
	return err
}

func (s *hangoutService) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) (*dto.PaginatedHangouts, error) {
	recordMetrics := s.metrics.StartRequest(ctx, "hangout", "list")

	hangouts, err := s.hangoutRepo.GetHangoutsByUserID(ctx, userID, pagination)
	if err != nil {
		recordMetrics("error")
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

	recordMetrics("success")
	return &dto.PaginatedHangouts{
		Data:       responseDTOs,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil

}
