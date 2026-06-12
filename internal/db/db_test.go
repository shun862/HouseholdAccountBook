package db

import (
	"database/sql"
	"testing"
)

func TestOpenDB(t *testing.T) {
	db, err := OpenDB(":memory:")
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}

func TestCreateUserTable(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}

	err = createUserTable(db)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// usersテーブルのカラム情報を取得
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		t.Fatalf("PRAGMA failed: %v", err)
	}
	defer rows.Close()

	cols := 0
	for rows.Next() {
		cols++
	}
	if cols != 3 {
		t.Errorf("users table is unusual. columns: %d", cols)
	}
}

func TestCreateExpenseTable(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}

	err = createExpenseTable(db)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// usersテーブルのカラム情報を取得
	rows, err := db.Query("PRAGMA table_info(expenses)")
	if err != nil {
		t.Fatalf("PRAGMA failed: %v", err)
	}
	defer rows.Close()

	cols := 0
	for rows.Next() {
		cols++
	}
	if cols != 6 {
		t.Errorf("expenses table is unusual. columns: %d", cols)
	}
}
