package db

import (
	"database/sql"

	"bitbucket.org/liamstask/goose/lib/goose"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/ziutek/mymysql/godrv"
)

type TxFunc func() error

type logger interface {
	Print(v ...interface{})
}

type Database interface {
	Close() error
	LogMode(enable bool)
	SetLogger(log logger)
	Migrate(dbName, dbPath, mirgrationsDir, migrationsEnv string) error
	Create(value interface{}) error
	CreateTx(value interface{}) error
	Read(out interface{}, query interface{}, args ...interface{}) error
	ReadAll(out interface{}) error
	Update(value interface{}, query interface{}, args ...interface{}) error
	UpdateTx(value interface{}, query interface{}, args ...interface{}) error
	Delete(value interface{}, query interface{}, args ...interface{}) error
	DeleteTx(value interface{}, query interface{}, args ...interface{}) error
	Save(value interface{}) error
	SaveTx(value interface{}) error
	Transaction(fn TxFunc) (err error)
}

type database struct {
	*gorm.DB
}

func New(name, path string, maxOpenConnections int) (Database, error) {
	db, err := gorm.Open(name, path)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(maxOpenConnections)
	return &database{db}, nil
}

func (db *database) db() *sql.DB {
	return db.DB.DB()
}

func (db *database) Migrate(dbName, dbPath, mirgrationsDir, migrationsEnv string) error {
	c := &goose.DBConf{
		MigrationsDir: mirgrationsDir,
		Env:           migrationsEnv,
		Driver:        driver(dbName, dbPath),
	}
	v, err := goose.GetMostRecentDBVersion(mirgrationsDir)
	if err != nil {
		return err
	}
	return goose.RunMigrationsOnDb(c, mirgrationsDir, v, db.db())
}

func (db *database) Close() error {
	return db.DB.Close()
}

func (db *database) Create(value interface{}) error {
	return db.DB.Create(value).Error
}

func (db *database) LogMode(enable bool) {
	db.DB.LogMode(enable)
}

func (db *database) SetLogger(log logger) {
	db.DB.SetLogger(log)
}

func (db *database) CreateTx(value interface{}) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Create(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) Read(out interface{}, query interface{}, args ...interface{}) error {
	return db.DB.Where(query, args).First(out).Error
}

func (db *database) ReadAll(out interface{}) error {
	return db.DB.Find(out).Error
}

func (db *database) Update(value interface{}, query interface{}, args ...interface{}) error {
	return db.DB.Where(query, args).Update(value).Error
}

func (db *database) UpdateTx(value interface{}, query interface{}, args ...interface{}) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Where(query, args).Update(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) Delete(value interface{}, query interface{}, args ...interface{}) error {
	return db.DB.Where(query, args).Delete(value).Error
}

func (db *database) DeleteTx(value interface{}, query interface{}, args ...interface{}) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Where(query, args).Delete(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) Save(value interface{}) error {
	return db.DB.Save(value).Error
}

func (db *database) SaveTx(value interface{}) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Save(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) Transaction(fn TxFunc) (err error) {
	tx, err := db.DB.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	return fn()
}

func driver(name, openStr string) goose.DBDriver {
	d := goose.DBDriver{
		Name:    name,
		OpenStr: openStr,
	}
	switch name {
	case "postgres":
		d.Import = "github.com/lib/pq"
		d.Dialect = &goose.PostgresDialect{}
	case "mymysql":
		d.Import = "github.com/ziutek/mymysql/godrv"
		d.Dialect = &goose.MySqlDialect{}
	case "mysql":
		d.Import = "github.com/go-sql-driver/mysql"
		d.Dialect = &goose.MySqlDialect{}
	case "sqlite3":
		d.Import = "github.com/mattn/go-sqlite3"
		d.Dialect = &goose.Sqlite3Dialect{}
	}
	return d
}
