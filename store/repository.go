package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/kwtryo/go-sample/clock"
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

	// sql.Openは接続確認が行われないため、ここで確認する。
	if err := db.Ping(); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	xdb := sqlx.NewDb(db, driver)
	return xdb, func() { _ = db.Close() }, nil
}

type Beginner interface {
	BeginTx(opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	Preparex(query string) (*sqlx.Stmt, error)
}

type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	QueryRow(query string, args ...any) *sqlx.Row
	Get(dest interface{}, query string, args ...any) error
	Select(dest interface{}, query string, args ...any) error
}

var (
	// インタフェースが期待通りに宣言されているか確認
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sql.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}
