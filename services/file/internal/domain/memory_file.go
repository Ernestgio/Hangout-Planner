package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemoryFile struct {
	ID            uuid.UUID `gorm:"primaryKey;type:char(36)"`
	OriginalName  string    `gorm:"type:varchar(255);not null"`
	FileExtension string    `gorm:"type:varchar(10);not null"`
	StoragePath   string    `gorm:"type:varchar(500);not null"`
	FileSize      int64     `gorm:"not null"`
	MimeType      string    `gorm:"type:varchar(100);not null"`
	FileStatus    string    `gorm:"type:varchar(50);not null"`
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	MemoryID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex"`
}

func (file *MemoryFile) BeforeCreate(tx *gorm.DB) (err error) {
	file.ID = uuid.New()
	return
}
