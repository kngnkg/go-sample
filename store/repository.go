package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/kwtryo/go-sample/config"
)

func New(cfg *config.Config) (*sqlx.DB, func(), error) {
	driver := "mysql"
	db, err := sql.Open(driver, fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	))
	if err != nil {
		return nil, nil, err
	}

	// Openは接続確認が行われないため、ここで確認する。
	if err := db.Ping(); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	xdb := sqlx.NewDb(db, driver)
	return xdb, func() { _ = db.Close() }, nil
}
