package repository

import (
	"Hertz-Scaffold/conf"
	"database/sql"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDbPool *sql.DB

func InitMysqlDb() {
	sqlDB, err := sql.Open("mysql", conf.AppConf.GetMysqlInfo().MysqlUrl)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(conf.AppConf.GetMysqlInfo().MaxIdleConn)                       // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(conf.AppConf.GetMysqlInfo().MaxOpenConn)                       // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Duration(conf.AppConf.GetMysqlInfo().MaxConnLifeTime)) //设置了连接可复用的最大时间
	SqlDbPool = sqlDB
}

func GetGormDb() (*gorm.DB, error) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: SqlDbPool}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
