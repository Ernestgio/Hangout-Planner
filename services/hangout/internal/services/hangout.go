package services

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HangoutService interface {
	UpdateHangout(id, userID uuid.UUID, req *dto.UpdateHangoutRequest) (*domain.Hangout, error)
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

func (s *hangoutService) UpdateHangout(id, userID uuid.UUID, req *dto.UpdateHangoutRequest) (*domain.Hangout, error) {
	var updatedHangout *domain.Hangout

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.hangoutRepo.WithTx(tx)

		existingHangout, err := txRepo.GetHangoutByID(id)
		if err != nil {
			return err
		}

		if *existingHangout.UserID != userID {
			return apperrors.ErrForbidden
		}

		if req.Title != "" {
			existingHangout.Title = req.Title
		}

		updatedHangout, err = txRepo.UpdateHangout(existingHangout)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedHangout, nil
}
