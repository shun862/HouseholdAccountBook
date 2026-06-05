package handler

import (
	"household_account_book/internal/consts"
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
)

type AddExpenseHandler struct {
	Repo  *repository.ExpenseRepository
	Store *sessions.CookieStore
}

type AddExpenseViewInfo struct {
	CurrentPage string
	Title       string
	TitleError  string
	ExpenseType int
	Amount      string
	AmountError string
	ExpenseDate string
	IsError     bool

	Message string
}

func NewAddExpenseHandler(repo *repository.ExpenseRepository,
	store *sessions.CookieStore) *AddExpenseHandler {
	return &AddExpenseHandler{Repo: repo, Store: store}
}

func (handler *AddExpenseHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		amountStr := r.FormValue("amount")
		expenseType, _ := strconv.Atoi(r.FormValue("type"))
		dateStr := r.FormValue("date")
		expenseDate := strings.ReplaceAll(dateStr, "-", "/")
		userId := GetUserID(r.Context())

		info := AddExpenseViewInfo{
			CurrentPage: consts.AddExpenseViewName,
			Title:       title,
			ExpenseType: expenseType,
			Amount:      amountStr,
			ExpenseDate: dateStr,
		}

		addExpenseValidation(&info)
		if info.IsError {
			temp, _ := template.ParseFiles(
				consts.AddExpenseFile,
				consts.HeaderFile,
			)
			temp.Execute(w, info)
			return
		}

		// DB追加
		amount, _ := strconv.Atoi(amountStr)
		if err := handler.Repo.InsertExpense(userId, expenseType, amount, title, expenseDate); err != nil {
			log.Fatal(err)
		}

		// フラッシュメッセージをセット
		session, _ := handler.Store.Get(r, consts.SessionName)
		session.AddFlash("登録しました")
		session.Save(r, w)

		// ★ PRGパターンでリダイレクト
		http.Redirect(w, r, consts.AddExpenseUrl, http.StatusSeeOther)
	} else {
		temp, err := template.ParseFiles(
			consts.AddExpenseFile,
			consts.HeaderFile,
		)
		if err != nil {
			log.Fatal(err)
		}
		info := AddExpenseViewInfo{
			CurrentPage: consts.AddExpenseViewName,
			ExpenseType: 0,
			ExpenseDate: time.Now().Format(consts.DateFormatH),
		}

		// セッション情報取得
		session, _ := handler.Store.Get(r, consts.SessionName)
		flashes := session.Flashes()
		session.Save(r, w)
		// フラッシュメッセージがあれば設定
		if len(flashes) > 0 {
			info.Message = flashes[0].(string)
		}
		temp.Execute(w, info)
	}
}

func addExpenseValidation(info *AddExpenseViewInfo) {
	title := info.Title
	amount, _ := strconv.Atoi(info.Amount)

	// タイトル
	if title == "" {
		info.TitleError = "タイトルを入力してください"
		info.IsError = true
	} else if len := len([]rune(title)); len > 20 {
		info.TitleError = "20文字以内で入力してください"
		info.IsError = true
	}

	// 金額
	if amount <= 0 {
		info.AmountError = "1以上の値を入力してください"
		info.IsError = true
	}
}
