package handler

import (
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"strings"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

type ViewInfo struct {
	Username      string
	UsernameError string
	PasswordError string
	HasError      bool
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (handler *UserHandler) RegisterHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		info := userRegisterValidation(username, password, handler.Repo)
		if info.HasError {
			temp, _ := template.ParseFiles(
				"../../web/templates/user_register.html",
			)
			temp.Execute(w, info)
			return
		}

		// DB登録
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err := handler.Repo.Insert(username, string(passwordHash)); err != nil {
			log.Fatal(err)
		}

		http.Redirect(
			w,
			r,
			"/login",
			http.StatusSeeOther,
		)
	} else {
		temp, err := template.ParseFiles("../../web/templates/user_register.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, ViewInfo{})
	}
}

func userRegisterValidation(username string, password string, repo *repository.UserRepository) ViewInfo {
	info := ViewInfo{
		Username: username,
	}

	// 共通バリデーション
	validation(username, password, &info)

	if info.HasError {
		return info
	}

	// 同一ユーザーの存在チェック
	hasSameUser, err := repo.CheckSameUser(username)
	if err != nil {
		log.Fatal(err)
	}
	if hasSameUser {
		info.UsernameError = "すでに登録済みのユーザーです"
		info.PasswordError = "すでに登録済みのユーザーです"
		info.HasError = true
	}

	return info
}

// 共通バリデーション
func validation(username string, password string, info *ViewInfo) {
	// ユーザー名
	if username == "" {
		info.UsernameError = "ユーザー名を入力してください"
		info.HasError = true
	} else if len := len([]rune(username)); len < 4 || len > 10 {
		info.UsernameError = "ユーザー名は4～10文字で入力してください"
		info.HasError = true
	}

	// パスワード
	if password == "" {
		info.PasswordError = "パスワードを入力してください"
		info.HasError = true
	} else if len := len([]rune(password)); len < 8 || len > 12 {
		info.PasswordError = "パスワードは8～12文字で入力してください"
		info.HasError = true
	}
}
