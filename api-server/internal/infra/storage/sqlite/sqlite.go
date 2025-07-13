package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/data/app.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
