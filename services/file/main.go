package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/app"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/constants/logmsg"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/logger"
)

func main() {
	ctx := context.Background()

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error(logmsg.ConfigLoadFailed, slog.Any("error", err))
		os.Exit(1)
	}

	// 2. Initialize logger
	logger.Init(cfg.Env, cfg.AppName)

	// 3. Create a new application instance
	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Error(ctx, logmsg.AppCreateFailed, err)
		os.Exit(1)
	}

	// 4. Start the application
	if err := application.Start(); err != nil {
		logger.Error(ctx, logmsg.AppTerminatedWithError, err)
		os.Exit(1)
	}

	logger.Info(ctx, logmsg.AppExitSuccess)
}
