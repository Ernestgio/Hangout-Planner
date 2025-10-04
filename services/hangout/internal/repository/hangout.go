package repository

import (
	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HangoutRepository interface {
	CreateHangout(hangout *domain.Hangout) (*domain.Hangout, error)
	GetHangoutByID(id uuid.UUID) (*domain.Hangout, error)
	UpdateHangout(hangout *domain.Hangout) (*domain.Hangout, error)
	DeleteHangout(id uuid.UUID) error
	WithTx(tx *gorm.DB) HangoutRepository
}

type hangoutRepository struct {
	db *gorm.DB
}

func NewHangoutRepository(db *gorm.DB) HangoutRepository {
	return &hangoutRepository{db: db}
}
func (r *hangoutRepository) CreateHangout(hangout *domain.Hangout) (*domain.Hangout, error) {
	if err := r.db.Create(hangout).Error; err != nil {
		return nil, err
	}
	return hangout, nil
}
func (r *hangoutRepository) GetHangoutByID(id uuid.UUID) (*domain.Hangout, error) {
	var hangout domain.Hangout
	if err := r.db.Preload("User").First(&hangout, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &hangout, nil
}

func (r *hangoutRepository) UpdateHangout(hangout *domain.Hangout) (*domain.Hangout, error) {
	if err := r.db.Save(hangout).Error; err != nil {
		return nil, err
	}
	return hangout, nil
}
func (r *hangoutRepository) DeleteHangout(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Hangout{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *hangoutRepository) WithTx(tx *gorm.DB) HangoutRepository {
	return &hangoutRepository{db: tx}
}
