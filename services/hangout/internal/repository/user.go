package repository

import (
	"context"
	"time"

	domain "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	CreateUser(context context.Context, user *domain.User) error
	GetUserByEmail(context context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db      *gorm.DB
	metrics *otel.MetricsRecorder
}

func NewUserRepository(db *gorm.DB, metrics *otel.MetricsRecorder) UserRepository {
	return &userRepository{db: db, metrics: metrics}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{db: tx, metrics: r.metrics}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(user).Error
	r.metrics.RecordDBOperation(ctx, "insert", "users", time.Since(start), 1)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	start := time.Now()
	var user domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	r.metrics.RecordDBOperation(ctx, "select", "users", time.Since(start), 1)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
