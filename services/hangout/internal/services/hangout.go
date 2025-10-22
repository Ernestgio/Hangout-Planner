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
	GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]*dto.HangoutListItemResponse, error)
}

type hangoutService struct {
	db          *gorm.DB
	hangoutRepo repository.HangoutRepository
}

func NewHangoutService(db *gorm.DB, hangoutRepo repository.HangoutRepository) HangoutService {
	return &hangoutService{
		db:          db,
		hangoutRepo: hangoutRepo,
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
		txRepo := s.hangoutRepo.WithTx(tx)

		existingHangout, err := txRepo.GetHangoutByID(ctx, id, userID)
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

		updatedHangout, err = txRepo.UpdateHangout(ctx, existingHangout)
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

func (s *hangoutService) GetHangoutsByUserID(ctx context.Context, userID uuid.UUID, pagination *dto.CursorPagination) ([]*dto.HangoutListItemResponse, error) {
	hangouts, err := s.hangoutRepo.GetHangoutsByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	response := mapper.HangoutsToListItemResponseDTOs(hangouts)
	return response, nil
}
