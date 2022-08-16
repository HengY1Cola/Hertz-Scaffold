package repository

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDbPool *sql.DB

func GetGormDb() (*gorm.DB, error) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: SqlDbPool}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
