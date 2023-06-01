package gormboot

import (
	"os"
	"strconv"
	"sync"

	"github.com/neoguojing/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db     *gorm.DB
	models = []interface{}{}
	once   sync.Once
)

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func init() {

}

func Init() {
	once.Do(func() {
		source := os.Getenv("db.source")
		dbPath := os.Getenv("db.path")
		maxIdle := os.Getenv("db.maxIdle")
		maxConn := os.Getenv("db.maxConn")
		enableLog := os.Getenv("db.log")

		if source == "" {
			panic("env db.source was empty")
		}

		if dbPath == "" {
			panic("env db.path was empty")
		}

		var err error
		db, err = gorm.Open(source, dbPath)
		if err != nil {
			panic(err)
		}

		if enableLog != "" {
			db.LogMode(true)
		}

		if maxConn != "" {
			maxC, err := strconv.Atoi(maxConn)
			if err == nil {
				db.DB().SetMaxOpenConns(maxC)
			} else {
				log.Error(err.Error())
			}

		}

		if maxIdle != "" {
			maxI, err := strconv.Atoi(maxIdle)
			if err == nil {
				db.DB().SetMaxIdleConns(maxI)
			} else {
				log.Error(err.Error())
			}

		}

		for _, m := range models {
			if !db.HasTable(m) {
				err = db.CreateTable(m).Error
				if err != nil {
					panic(err)
				}
			}
		}

		if len(models) != 0 {
			err = db.AutoMigrate(models...).Error
			if err != nil {
				panic(err)
			}
		}
	})
}

func getDB() *gorm.DB {
	err := db.DB().Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GetDB() *gorm.DB {
	return getDB()
}

func Destroy() {
	db.Close()
}
