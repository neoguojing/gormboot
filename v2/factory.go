package gormboot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/neoguojing/log"

	"gorm.io/gorm"
)

type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	SQLite     DatabaseType = "sqlite"
	ClickHouse DatabaseType = "clickhouse"
)

var (
	DefaultDB *Factory
)

func init() {
	basepath := os.Getenv("DB_PATH")
	dbPath := filepath.Join(basepath, "sqlite3.db")
	DefaultDB = New(DefaultSqliteConfig(dbPath))
}

func DefaulMysqlConfig(user, password, host, database string, port int) *DBConfig {
	cfg := &DBConfig{
		Source:   MySQL,
		User:     user,
		Password: password,
		Host:     host,
		DataBase: database,
		Port:     port,
	}
	return cfg
}

func DefaulClickhouseConfig(user, password, host, database string, port int) *DBConfig {
	cfg := &DBConfig{
		Source:   MySQL,
		User:     user,
		Password: password,
		Host:     host,
		DataBase: database,
		Port:     port,
	}
	return cfg
}

func DefaultSqliteConfig(filePath string) *DBConfig {
	cfg := &DBConfig{
		Source:   SQLite,
		FilePath: filePath,
	}
	return cfg
}

type DBConfig struct {
	Source       DatabaseType
	User         string
	Password     string
	DataBase     string
	FilePath     string
	Host         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	MaxIdle      int
	MaxOpen      int
}

func checkEnv(source DatabaseType) error {
	switch source {
	case MySQL, ClickHouse:
		if os.Getenv("db.host") == "" {
			return errors.New("db.host is not set")
		}
		if os.Getenv("db.port") == "" {
			return errors.New("db.port is not set")
		}
		if os.Getenv("db.user") == "" {
			return errors.New("db.user is not set")
		}
		if os.Getenv("db.password") == "" {
			return errors.New("db.password is not set")
		}
		if os.Getenv("db.database") == "" {
			return errors.New("db.database is not set")
		}
	case SQLite:
		if os.Getenv("db.filepath") == "" {
			return errors.New("db.filepath is not set")
		}
	default:
		return errors.New("unsupported database type: " + string(source))
	}
	return nil
}

func BuildByEnv(source DatabaseType) *DBConfig {
	if err := checkEnv(source); err != nil {
		panic(fmt.Sprintf("failed to check env: %v", err))
	}
	dbConfig := DBConfig{}
	portInt, err := strconv.Atoi(os.Getenv("db.port"))
	if err != nil {
		panic(fmt.Sprintf("failed to convert port to int: %v", err))
	}

	switch source {
	case MySQL:
		dbConfig = DBConfig{
			Host:     os.Getenv("db.host"),
			Port:     portInt,
			User:     os.Getenv("db.user"),
			Password: os.Getenv("db.password"),
			DataBase: os.Getenv("db.database"),
		}
	case SQLite:
		dbConfig = DBConfig{
			FilePath: os.Getenv("db.filepath"),
		}
	case ClickHouse:
		dbConfig = DBConfig{
			Host:     os.Getenv("db.host"),
			Port:     portInt,
			User:     os.Getenv("db.user"),
			Password: os.Getenv("db.password"),
			DataBase: os.Getenv("db.database"),
		}
	default:
		panic(fmt.Sprintf("unsupported database type: %s", source))
	}
	return &dbConfig
}

func (d *DBConfig) Pool(maxConn uint, maxIdle uint) *DBConfig {
	d.MaxIdle = int(maxIdle)
	d.MaxIdle = int(maxConn)
	return d
}

func (cfg *DBConfig) DSN() string {

	switch cfg.Source {
	case MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DataBase)
	case SQLite:
		return cfg.FilePath
	case ClickHouse:
		return fmt.Sprintf("tcp://%s:%s@%s:%d?database=%s", cfg.User, cfg.Password,
			cfg.Host, cfg.Port, cfg.DataBase)

	default:
		panic(fmt.Sprintf("unsupported database type: %s", cfg.Source))
		return ""
	}
}

type Factory struct {
	models []interface{}
	db     *gorm.DB
}

func NewByEnv(dbType DatabaseType) *Factory {
	cfg := BuildByEnv(dbType)

	return New(cfg)
}

func New(cfg *DBConfig) *Factory {
	var db *gorm.DB
	switch cfg.Source {
	case MySQL:
		db = NewMysql(cfg)
	case SQLite:
		db = NewSqlite(cfg)
	case ClickHouse:
		db = NewClickhouse(cfg)
	default:
		panic("invalid db source")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}

	return &Factory{
		db:     db,
		models: make([]interface{}, 0),
	}
}

func (d *Factory) RegisterModel(model ...interface{}) {
	d.models = append(d.models, model...)
}

func (d *Factory) AutoMigrate() *Factory {

	for _, m := range d.models {
		if !d.db.Migrator().HasTable(m) {
			err := d.db.Migrator().CreateTable(m)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
		}
	}

	if len(d.models) != 0 {
		err := d.db.AutoMigrate(d.models...)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}

	return d
}

func (d *Factory) Close() error {
	if d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		log.Error(err.Error())
	}
	return sqlDB.Close()
}

func (d *Factory) DB() *gorm.DB {
	return d.db
}
