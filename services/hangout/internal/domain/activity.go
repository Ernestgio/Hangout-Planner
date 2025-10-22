package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID *uuid.UUID `gorm:"type:char(36)"`
	User   User       `gorm:"foreignKey:UserID"`

	Hangouts []*Hangout `gorm:"many2many:hangout_activities;"`
}

func (activity *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	activity.ID = uuid.New()
	return
}
