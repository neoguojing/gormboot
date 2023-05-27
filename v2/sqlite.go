package gormboot

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(cfg *DBConfig) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
