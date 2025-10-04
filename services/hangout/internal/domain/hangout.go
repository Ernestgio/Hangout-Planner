package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hangout struct {
	ID          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
