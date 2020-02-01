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

func New(name, path string, maxOpenConnections int) (db *gorm.DB, err error) {
	if db, err = gorm.Open(name, path); err != nil {
		return
	}
	db.DB().SetMaxOpenConns(maxOpenConnections)
	return
}

func Migrate(db *sql.DB, dbName, dbPath, mirgrationsDir, migrationsEnv string) error {
	c := &goose.DBConf{
		MigrationsDir: mirgrationsDir,
		Env:           migrationsEnv,
		Driver:        driver(dbName, dbPath),
	}
	v, err := goose.GetMostRecentDBVersion(mirgrationsDir)
	if err != nil {
		return err
	}
	return goose.RunMigrationsOnDb(c, mirgrationsDir, v, db)
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
