package handler

import (
	"household_account_book/internal/consts"
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"text/template"
)

type AddExpenseHandler struct {
	Repo *repository.ExpenseRepository
}

type AddExpenseViewInfo struct {
	CurrentPage string
	Title       string
	TitleError  string
	Amount      string
	AmountError string
	IsError     bool
}

func NewAddExpenseHandler(repo *repository.ExpenseRepository) *AddExpenseHandler {
	return &AddExpenseHandler{Repo: repo}
}

func (handler *AddExpenseHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		temp, err := template.ParseFiles(
			consts.AddExpenseFile,
			consts.HeaderFile,
		)
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, AddExpenseViewInfo{CurrentPage: consts.AddExpenseViewName})
	}
}
