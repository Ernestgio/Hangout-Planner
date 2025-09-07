package db

import (
	"fmt"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/models"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, func() error, error) {
	dsn := buildDSN(cfg)

	gormDB, err := OpenGORM(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("open gorm: %w", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get sql.DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping db: %w", err)
	}

	closer := func() error {
		return db.Close()
	}

	return gormDB, closer, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
