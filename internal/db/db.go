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

	// 外部キー制約を有効化
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return db, err
	}

	// usersテーブル作成
	err = createUserTable(db)
	if err != nil {
		return db, err
	}

	// expensesテーブル作成
	err = createExpenseTable(db)
	return db, err
}

func createUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
    `
	_, err := db.Exec(query)
	return err
}

func createExpenseTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		expense_type INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK(amount > 0),
		title TEXT NOT NULL,
		expense_date TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
    `
	_, err := db.Exec(query)
	return err
}
