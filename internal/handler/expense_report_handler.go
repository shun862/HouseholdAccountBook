package handler

import (
	"encoding/json"
	"household_account_book/internal/consts"
	"household_account_book/internal/repository"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type ExpenseReportHandler struct {
	Repo *repository.ExpenseRepository
}

type ExpenseReportViewInfo struct {
	CurrentPage string
	DisplayType int
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
		displayType := r.FormValue("display")
		userId := GetUserID(r.Context())
		info := createViewInfo(userId, displayType, handler.Repo)
		temp.Execute(w, info)
	}
}

func createViewInfo(userId int, displayTypeStr string, repo *repository.ExpenseRepository) ExpenseReportViewInfo {
	// 対象月数
	var len int
	// 表示期間タイプを取得
	displayType, err := strconv.Atoi(displayTypeStr)
	if err != nil {
		displayType = 0
	}
	if displayType == 1 {
		len = 3
	} else {
		len = 6
	}

	now := time.Now()
	// 対象日付の収支情報を取得
	startDate := getMonthStart(now.AddDate(0, -len+1, 0))
	endDate := getMonthEnd(now)
	expenses, _ := repo.FindExpenses(userId, startDate, endDate)

	// 月ごとの集計用マップ
	type summary struct {
		income  int
		expense int
	}
	monthly := make(map[string]*summary, len)
	for i := 0; i < len; i++ {
		key := startDate.AddDate(0, i, 0).Format(consts.DateFormatJ)
		monthly[key] = &summary{}
	}

	for _, e := range expenses {
		// 日付を time.Time に変換
		t, err := time.Parse(consts.DateFormatS, e.ExpenseDate)
		if err != nil {
			continue
		}

		// 6か月以内のデータだけ対象
		if t.Before(startDate) || t.After(endDate) {
			continue
		}

		// キーは "2026-06" のようにする
		key := t.Format(consts.DateFormatJ)

		// ExpenseType で振り分け
		if e.ExpenseType == 1 {
			monthly[key].income += e.Amount
		} else {
			monthly[key].expense += e.Amount
		}
	}

	// MonthlyReport の配列に変換
	var reports []MonthlyReport
	for key, s := range monthly {
		reports = append(reports, MonthlyReport{
			Month:   key,
			Income:  s.income,
			Expense: s.expense,
		})
	}

	// mapのrangeはランダムのため、昇順にソートする
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Month < reports[j].Month
	})

	// Chart.jsで使用するためのjsonデータに変換
	var labels []string
	var totalIncomes []int
	var totalExpenses []int
	for _, report := range reports {
		labels = append(labels, report.Month)
		totalIncomes = append(totalIncomes, report.Income)
		totalExpenses = append(totalExpenses, report.Expense)
	}

	labelsJSON, _ := json.Marshal(labels)
	incomeJSON, _ := json.Marshal(totalIncomes)
	expenseJSON, _ := json.Marshal(totalExpenses)

	info := ExpenseReportViewInfo{
		CurrentPage: consts.ExpenseReportViewName,
		DisplayType: displayType,
		Reports:     reports,
		MonthJSON:   template.JS(labelsJSON),
		IncomeJSON:  template.JS(incomeJSON),
		ExpenseJSON: template.JS(expenseJSON),
	}
	return info
}

func getMonthEnd(t time.Time) time.Time {
	// 翌月1日の0時を取得
	d := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	// 1日引いて月末日を取得
	return d.AddDate(0, 0, -1)
}

func getMonthStart(t time.Time) time.Time {
	// 翌月1日の0時を取得
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
