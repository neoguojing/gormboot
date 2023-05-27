package gormboot

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func NewClickhouse(cfg *DBConfig) *gorm.DB {
	conn := clickhouse.Open(cfg.DSN())
	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

