package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hangout struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:char(36)"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Date        *time.Time `gorm:"not null" json:"date"`
	Status      *string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	UserID *uuid.UUID `gorm:"type:char(36)"`
	User   User       `gorm:"foreignKey:UserID"`
}

func (hangout *Hangout) BeforeCreate(tx *gorm.DB) (err error) {
	hangout.ID = uuid.New()
	return
}
