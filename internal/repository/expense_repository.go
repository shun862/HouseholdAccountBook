package repository

import (
	"database/sql"
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
