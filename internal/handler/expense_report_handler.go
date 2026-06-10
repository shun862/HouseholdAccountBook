package handler

import (
	"encoding/json"
	"household_account_book/internal/consts"
	"household_account_book/internal/repository"
	"html/template"
	"log"
	"net/http"
)

type ExpenseReportHandler struct {
	Repo *repository.ExpenseRepository
}

type ExpenseReportViewInfo struct {
	CurrentPage string
	Reports     []MonthlyReport
	MonthJSON   template.JS
	IncomeJSON  template.JS
	ExpenseJSON template.JS
}

type MonthlyReport struct {
	Month   string
	Income  int
	Expense int
}

func NewExpenseReportHandler(repo *repository.ExpenseRepository) *ExpenseReportHandler {
	return &ExpenseReportHandler{Repo: repo}
}

func (handler *ExpenseReportHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles(
			consts.ExpenseReportFile,
			consts.HeaderFile,
		)
		if err != nil {
			log.Fatal(err)
		}
		// startDate := r.FormValue("")
		// userId := GetUserID(r.Context())
		// 検索用のフォーマットに変換
		// start, _ := time.Parse(consts.DateFormatH, startDate)
		// end, _ := time.Parse(consts.DateFormatH, endDate)
		// expenses, _ := handler.Repo.FindExpenses(userId, start, end)
		info := createViewInfo()
		temp.Execute(w, info)
	}
}

func createViewInfo() ExpenseReportViewInfo {
	reports := []MonthlyReport{
		{"2024年12月", 125000, 80000},
		{"2025年1月", 150000, 90000},
		{"2025年2月", 190000, 120000},
		{"2025年3月", 160000, 100000},
		{"2025年4月", 200000, 130000},
		{"2025年5月", 250000, 153450},
	}

	var labels []string
	var incomes []int
	var expenses []int

	for _, report := range reports {
		labels = append(labels, report.Month)
		incomes = append(incomes, report.Income)
		expenses = append(expenses, report.Expense)
	}

	labelsJSON, _ := json.Marshal(labels)
	incomeJSON, _ := json.Marshal(incomes)
	expenseJSON, _ := json.Marshal(expenses)

	info := ExpenseReportViewInfo{
		CurrentPage: consts.ExpenseReportViewName,
		Reports:     reports,
		MonthJSON:   template.JS(labelsJSON),
		IncomeJSON:  template.JS(incomeJSON),
		ExpenseJSON: template.JS(expenseJSON),
	}
	return info
}
