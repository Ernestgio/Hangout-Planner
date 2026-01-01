package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Memory struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_hangout_name,priority:2"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	HangoutID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_hangout_name,priority:1"`
	Hangout   Hangout   `gorm:"foreignKey:HangoutID"`

	UserID uuid.UUID `gorm:"type:char(36);not null"`
	User   User      `gorm:"foreignKey:UserID"`
}

func (memory *Memory) BeforeCreate(tx *gorm.DB) (err error) {
	memory.ID = uuid.New()
	return
}
