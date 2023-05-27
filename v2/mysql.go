package gormboot

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(cfg *DBConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
