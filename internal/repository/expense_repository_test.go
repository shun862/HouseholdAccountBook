package repository

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestExpenseTable(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS expenses(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			expense_type INTEGER NOT NULL,
			amount INTEGER NOT NULL CHECK(amount > 0),
			title TEXT NOT NULL,
			expense_date TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
    `)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestInsertExpense(t *testing.T) {
	db := setupTestExpenseTable(t)
	defer db.Close()

	repo := NewExpenseRepository(db)
	err := repo.InsertExpense(1, 0, 200, "ランチ", "2026/06/12")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	var count int
	err = db.QueryRow(`
        SELECT Count(*)
        FROM expenses
        WHERE user_id = ?
          AND expense_type = ?
          AND amount = ?
          AND title = ?
          AND expense_date = ?
    `,
		1,
		0,
		200,
		"ランチ",
		"2026/06/12",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to query inserted data: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestFindExpenses(t *testing.T) {
	db := setupTestExpenseTable(t)
	defer db.Close()

	repo := NewExpenseRepository(db)
	testData := []struct {
		userId      int
		expenseType int
		amount      int
		title       string
		date        string
	}{
		{1, 0, 100, "A", "2026/06/01"},
		{1, 1, 200, "B", "2026/05/01"},
		{1, 0, 300, "C", "2026/06/03"},
	}

	for _, e := range testData {
		err := repo.InsertExpense(e.userId, e.expenseType, e.amount, e.title, e.date)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	expenses, err := repo.FindExpenses(1, start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(expenses) != 2 {
		t.Fatalf("expected 2 expenses, got %d", len(expenses))
	}

	// 日付降順になっていること
	if expenses[0].ExpenseDate != "2026/06/03" {
		t.Errorf(
			"expected first date 2026/06/03, got %s",
			expenses[0].ExpenseDate,
		)
	}

	if expenses[1].ExpenseDate != "2026/06/01" {
		t.Errorf(
			"expected second date 2026/06/01, got %s",
			expenses[1].ExpenseDate,
		)
	}
}

func TestFindExpenses_Empty(t *testing.T) {
	db := setupTestExpenseTable(t)
	defer db.Close()
	repo := NewExpenseRepository(db)

	testData := []struct {
		userId      int
		expenseType int
		amount      int
		title       string
		date        string
	}{
		{1, 0, 100, "A", "2026/04/01"},
		{1, 1, 200, "B", "2026/05/01"},
		{1, 0, 300, "C", "2026/07/03"},
	}

	for _, e := range testData {
		err := repo.InsertExpense(e.userId, e.expenseType, e.amount, e.title, e.date)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	result, err := repo.FindExpenses(1, start, end)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
