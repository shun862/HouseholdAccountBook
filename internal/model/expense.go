package model

type Expense struct {
	Id          int
	UserId      int
	ExpenseType int
	Amount      int
	Title       string
	ExpenseDate string
}
