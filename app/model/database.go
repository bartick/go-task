package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DBTX is an abstraction over sqlx.DB and sqlx.Tx.
type DBTX interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func InitDatabases(config DatabaseConfig) (DBTX, error) {
	var dsn string
	if config.DBUser == "" && config.DBPass == "" {
		dsn = fmt.Sprintf("tcp(%s:%s)/%s?parseTime=true",
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.DBUser,
			config.DBPass,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
