package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func New(driver, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
