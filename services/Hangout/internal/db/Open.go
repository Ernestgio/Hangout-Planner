package db

import (
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenGORM(dsn string) (*gorm.DB, error) {
	return gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
}
