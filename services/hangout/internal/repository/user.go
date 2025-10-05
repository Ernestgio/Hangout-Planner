package repository

import (
	"context"

	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	CreateUser(context context.Context, user *domain.User) error
	GetUserByEmail(context context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
