package repository

import (
	"database/sql"
	"household_account_book/internal/consts"
	"household_account_book/internal/model"
	"sort"
	"time"
)

type ExpenseRepository struct {
	DB *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (repo *ExpenseRepository) InsertExpense(userId int, expenseType int,
	amount int, title string, expenseDate string) error {
	query := `
	INSERT INTO expenses (user_id, expense_type, amount, title, expense_date)
	VALUES(?, ?, ?, ?, ?)
	`
	_, err := repo.DB.Exec(query, userId, expenseType, amount, title, expenseDate)
	return err
}

func (repo *ExpenseRepository) FindExpenses(userId int, start time.Time, end time.Time) ([]model.Expense, error) {
	// 件数を先に取得
	var count int
	repo.DB.QueryRow("SELECT COUNT(*) FROM expenses WHERE user_id = ?", userId).Scan(&count)

	// レコード取得
	expenses := make([]model.Expense, 0, count)
	rows, err := repo.DB.Query("SELECT * FROM expenses WHERE user_id = ?", userId)
	if err != nil {
		return expenses, err
	}

	for rows.Next() {
		var e model.Expense
		if err = rows.Scan(&e.Id, &e.UserId, &e.ExpenseType,
			&e.Amount, &e.Title, &e.ExpenseDate); err != nil {
			return expenses, err
		}
		expenseDate, _ := time.Parse(consts.DateFormatS, e.ExpenseDate)
		// 対象日付内なら追加
		if !expenseDate.After(end) && !expenseDate.Before(start) {
			expenses = append(expenses, e)
		}
	}
	// 日付降順でソート
	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].ExpenseDate > expenses[j].ExpenseDate
	})
	return expenses, err
}
