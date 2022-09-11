package db

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/pressly/goose"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrUnknownDialect = errors.New("Unknown database dialect")
)

type Database interface {
	Close() error

	Create(value interface{}) error

	CreateTx(value interface{}) error

	Delete(value interface{}, query interface{}, args ...interface{}) error

	DeleteTx(value interface{}, query interface{},
		args ...interface{}) error

	EnableLogging(enable bool)

	Migrate(dir string) error

	Open(path string) error

	Read(out interface{}, query interface{}, args ...interface{}) error

	ReadAll(out interface{}) error

	Save(value interface{}) error

	SaveTx(value interface{}) error

	SetLogger(log Logger)

	SetMaxOpenConnections(max int)

	Tx(fn TxFn) (err error)

	Update(value interface{}, query interface{}, args ...interface{}) error

	UpdateTx(value interface{}, query interface{},
		args ...interface{}) error
}

type DialectFn func(dsn string) gorm.Dialector

type Logger logger.Writer

type TxFn func(tx *gorm.DB) error

type database struct {
	*gorm.DB
	dialect            string
	logger             Logger
	loggingEnabled     bool
	maxOpenConnections int
}

func New(dialect string) Database {
	return &database{dialect: dialect}
}

func (db *database) Close() error {
	sqlDb, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func (db *database) Create(value interface{}) error {
	return db.DB.Create(value).Error
}

func (db *database) CreateTx(value interface{}) error {
	return db.Tx(func(tx *gorm.DB) error {
		return tx.Create(value).Error
	})
}

func (db *database) Delete(value interface{}, query interface{},
	args ...interface{}) error {
	return db.DB.Where(query, args...).Delete(value).Error
}

func (db *database) DeleteTx(value interface{}, query interface{},
	args ...interface{}) error {
	return db.Tx(func(tx *gorm.DB) error {
		return tx.Where(query, args...).Delete(value).Error
	})
}

func (db *database) EnableLogging(enable bool) {
	db.loggingEnabled = enable
}

func (db *database) Migrate(dir string) error {
	sqlDb, err := db.DB.DB()
	if err != nil {
		return err
	}
	if err := goose.SetDialect(db.dialect); err != nil {
		return err
	}
	return goose.Up(sqlDb, dir)
}

func (db *database) Open(path string) error {
	dbCfg := getDbConfig(db)
	dialectFn, err := getGormDialectorFn(db.dialect)
	if err != nil {
		return err
	}
	db.DB, err = gorm.Open(dialectFn(path), &dbCfg)
	if err != nil {
		return err
	}
	if sqlDb, err := db.DB.DB(); err == nil {
		sqlDb.SetMaxOpenConns(db.maxOpenConnections)
	} else {
		return err
	}
	return nil
}

func (db *database) Read(out interface{}, query interface{},
	args ...interface{}) error {
	return db.DB.Where(query, args...).First(out).Error
}

func (db *database) ReadAll(out interface{}) error {
	return db.DB.Find(out).Error
}

func (db *database) Save(value interface{}) error {
	return db.DB.Save(value).Error
}

func (db *database) SaveTx(value interface{}) error {
	return db.Tx(func(tx *gorm.DB) error {
		return tx.Save(value).Error
	})
}

func (db *database) SetLogger(log Logger) {
	db.logger = log
}

func (db *database) SetMaxOpenConnections(max int) {
	db.maxOpenConnections = max
}

func (db *database) Tx(fn TxFn) (err error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit().Error
	}()
	return fn(tx)
}

func (db *database) Update(value interface{}, query interface{},
	args ...interface{}) error {
	return db.DB.Where(query, args...).Updates(value).Error
}

func (db *database) UpdateTx(value interface{}, query interface{},
	args ...interface{}) error {
	return db.Tx(func(tx *gorm.DB) error {
		return tx.Where(query, args...).Updates(value).Error
	})
}

func getDbConfig(db *database) gorm.Config {
	var cfg gorm.Config
	loggerCfg := logger.Config{
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Silent,
		SlowThreshold:             200 * time.Millisecond,
	}
	if db.loggingEnabled {
		loggerCfg.LogLevel = logger.Warn
	}
	if db.logger != nil {
		cfg.Logger = logger.New(db.logger, loggerCfg)
	} else {
		cfg.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			loggerCfg,
		)
	}
	return cfg
}

func getGormDialectorFn(dialect string) (DialectFn, error) {
	switch dialect {
	case "mysql":
		return mysql.Open, nil
	case "postgres":
		return postgres.Open, nil
	case "sqlite3":
		return sqlite.Open, nil
	case "sqlserver":
		return sqlserver.Open, nil
	}
	return nil, ErrUnknownDialect
}
