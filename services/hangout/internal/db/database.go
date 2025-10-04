package db

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.DBConfig) (*gorm.DB, func() error, error) {
	dsn := buildDSN(cfg)
	gormDB, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	closer := func() error {
		return db.Close()
	}

	return gormDB, closer, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.Hangout{})
	if err != nil {
		return err
	}

	return nil
}
