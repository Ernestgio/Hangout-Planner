package db

import (
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/config"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func buildDSN(cfg *config.Config) string {
	mdcfg := mysqlDriver.Config{
		User:      cfg.DBUser,
		Passwd:    cfg.DBPassword,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		DBName:    cfg.DBName,
		ParseTime: true,
		Loc:       time.UTC,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	return mdcfg.FormatDSN()
}
