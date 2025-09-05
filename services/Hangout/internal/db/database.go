package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, func() error, error) {
	dsn := buildDSN(cfg)

	gormDB, err := OpenGORM(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("open gorm: %w", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get sql.DB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping db: %w", err)
	}

	closer := func() error {
		return db.Close()
	}

	return gormDB, closer, nil
}
