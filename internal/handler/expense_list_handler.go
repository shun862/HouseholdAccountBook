package handler

import (
	"household_account_book/internal/consts"
	"household_account_book/internal/model"
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"text/template"
	"time"
)

type ExpenseListHandler struct {
	Repo *repository.ExpenseRepository
}

type ExpenseListViewInfo struct {
	CurrentPage string
	StartDate   string
	EndDate     string
	Expenses    []model.Expense
}

func NewExpenseListHandler(repo *repository.ExpenseRepository) *ExpenseListHandler {
	return &ExpenseListHandler{Repo: repo}
}

func (handler *ExpenseListHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles(
			consts.ExpenseListFile,
			consts.HeaderFile,
		)
		if err != nil {
			log.Fatal(err)
		}
		startDate := r.FormValue("start")
		endDate := r.FormValue("end")
		if startDate == "" || endDate == "" {
			now := time.Now()
			// 1か月前の日付
			startDate = now.AddDate(0, -1, 0).Format(consts.DateFormatH)
			// 現在日付
			endDate = now.Format(consts.DateFormatH)
		}
		userId := GetUserID(r.Context())
		// 検索用のフォーマットに変換
		start, _ := time.Parse(consts.DateFormatH, startDate)
		end, _ := time.Parse(consts.DateFormatH, endDate)
		expenses, _ := handler.Repo.FindExpenses(userId, start, end)
		info := ExpenseListViewInfo{
			CurrentPage: consts.ExpenseListViewName,
			Expenses:    expenses,
			StartDate:   startDate,
			EndDate:     endDate,
		}
		temp.Execute(w, info)
	}
}
