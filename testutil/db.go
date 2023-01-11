package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func OpenDbForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	address := "docker.for.mac.localhost"
	port := 33306
	if _, defined := os.LookupEnv("CI"); defined {
		address = "127.0.0.1"
		port = 3306
	}

	driver := "mysql"
	db, err := sql.Open(
		driver,
		fmt.Sprintf("todo:todo@tcp(%s:%d)/todo?parseTime=true", address, port),
	)
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(
		func() { _ = db.Close() },
	)
	return sqlx.NewDb(db, driver)
}
