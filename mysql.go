package gormboot

import (
	"log"
	"os"
	"strconv"
	"sync"

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
			log.Fatal("env db.source was empty")
		}

		if dbPath == "" {
			log.Fatal("env db.path was empty")
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
				log.Print(err.Error())
			}

		}

		if maxIdle != "" {
			maxI, err := strconv.Atoi(maxIdle)
			if err == nil {
				db.DB().SetMaxIdleConns(maxI)
			} else {
				log.Print(err.Error())
			}

		}

		for _, m := range models {
			if !db.HasTable(m) {
				err = db.CreateTable().Error
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if len(models) != 0 {
			err = db.AutoMigrate(models...).Error
			if err != nil {
				log.Fatal(err)
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
