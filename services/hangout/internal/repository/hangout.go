package repository

import (
	"fmt"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HangoutRepository interface {
	WithTx(tx *gorm.DB) HangoutRepository
	CreateHangout(hangout *domain.Hangout) (*domain.Hangout, error)
	GetHangoutByID(id uuid.UUID) (*domain.Hangout, error)
	UpdateHangout(hangout *domain.Hangout) (*domain.Hangout, error)
	DeleteHangout(id uuid.UUID) error
	GetHangoutsByUserID(userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error)
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

func (r *hangoutRepository) GetHangoutsByUserID(userID uuid.UUID, pagination *dto.CursorPagination) ([]domain.Hangout, error) {
	var hangouts []domain.Hangout

	query := r.db.Model(&domain.Hangout{}).Where("user_id = ?", userID)

	if pagination.AfterID != nil {
		sortByColumn := pagination.GetSortBy()
		subQuery := r.db.Model(&domain.Hangout{}).Select(sortByColumn).Where("id = ?", *pagination.AfterID)

		comparisonOperator := ">="
		if pagination.GetSortDir() == constants.SortDirectionDesc {
			comparisonOperator = "<="
		}

		query = query.Where(fmt.Sprintf("%s %s (?)", sortByColumn, comparisonOperator), subQuery)

		idSubQuery := r.db.Model(&domain.Hangout{}).Select("id").Where(sortByColumn+" = (?)", subQuery)
		query = query.Where("id NOT IN (?)", idSubQuery)
	}

	err := query.
		Order(pagination.GetOrderByClause()).
		Limit(pagination.GetLimit()).
		Find(&hangouts).Error

	if err != nil {
		return nil, err
	}

	return hangouts, nil
}
