package database

import (
	"context"
	"fmt"
	"time"

	"github.com/abibby/eztvrss/config"
	"github.com/abibby/eztvrss/migrations"
	mysqlDialect "github.com/abibby/salusa/database/dialects/mysql"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Init(ctx context.Context) error {
	mysqlDialect.UseMySql()

	cfg := mysql.Config{
		User:   config.DBUsername,
		Passwd: config.DBPassword,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", config.DBHost, config.DBPort),
		DBName: config.DBDatabase,
	}
	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = migrations.Use().Up(ctx, db)
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func WithDB() router.MiddlewareFunc {
	return request.WithDB(DB)
}

func Tx(ctx context.Context, cb func(*sqlx.Tx) error) error {
	tx, err := DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = cb(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
