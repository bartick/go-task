package database

import (
	"fmt"
	"time"

	"github.com/bartick/ringover-task/app/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDatabases(config model.DatabaseConfig) (*sqlx.DB, error) {
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
