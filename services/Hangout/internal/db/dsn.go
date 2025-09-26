package db

import (
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func buildDSN(cfg *config.Config) string {
	mdcfg := mysqlDriver.Config{
		User:      cfg.DBUser,
		Passwd:    cfg.DBPassword,
		Addr:      fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		DBName:    cfg.DBName,
		Net:       constants.DefaultDBNetwork,
		ParseTime: true,
		Loc:       time.UTC,
		Params: map[string]string{
			"charset": constants.DefaultDBCharset,
		},
	}
	return mdcfg.FormatDSN()
}
