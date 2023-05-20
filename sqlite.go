package gormboot

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var instance *gorm.DB
var onceInstance sync.Once

func InitSqlite(dbPath string) (*gorm.DB, error) {
	onceInstance.Do(func() {
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic(err)
		}
		instance = db
	})
	return instance, nil
}

func GetSqlite() *gorm.DB {
	return instance
}
