package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	// DB取得
	db, err := sql.Open("sqlite3", "household_account_book.db")
	if err != nil {
		return nil, err
	}

	// usersテーブル作成
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT NOT NULL,
		password TEXT NOT NULL
	)
    `
	_, err = db.Exec(query)

	return db, err
}
