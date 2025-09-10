package db

import (
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/constants"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func buildDSN(cfg *config.Config) string {
	mdcfg := mysqlDriver.Config{
		User:      cfg.DBUser,
		Passwd:    cfg.DBPassword,
		Addr:      fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		DBName:    cfg.DBName,
		ParseTime: true,
		Loc:       time.UTC,
		Params: map[string]string{
			"charset": constants.DefaultDBCharset,
		},
	}
	return mdcfg.FormatDSN()
}
