package handler

import (
	"household_account_book/internal/consts"
	"household_account_book/internal/model"
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
			showView(w, consts.UserRegisterFile, info)
			return
		}

		// DB登録
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err := handler.Repo.Insert(username, string(passwordHash)); err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, consts.LoginUrl, http.StatusSeeOther)
	} else {
		showView(w, consts.UserRegisterFile, ViewInfo{})
	}
}

func (handler *UserHandler) LoginHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		info := ViewInfo{
			Username: username,
		}
		validation(username, password, &info)
		if info.HasError {
			showView(w, consts.LoginFile, info)
			return
		}

		// ユーザー情報取得
		user, err := handler.Repo.FindUser(username, password)
		if err != nil {
			log.Fatal(err)
		}
		// 未登録ユーザーの場合
		if user == nil {
			info.UsernameError = "未登録のユーザーです"
			info.PasswordError = "未登録のユーザーです"
			info.HasError = true
			showView(w, consts.LoginFile, info)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 1)
		claims := &model.JwtCustomClaims{
			Id:       user.Id,
			Username: username,
			Password: user.Password,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		// トークン生成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// トークンに署名
		tokenString, _ := token.SignedString([]byte(consts.JwtKey))

		// Cookie に JWT を保存
		http.SetCookie(w, &http.Cookie{
			Name:     consts.TokenName,
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
			Path:     consts.CookiePath,
		})

		http.Redirect(w, r, consts.AddExpenseUrl, http.StatusSeeOther)
	} else {
		showView(w, consts.LoginFile, ViewInfo{})
	}
}

func showView(w http.ResponseWriter, fileName string, info ViewInfo) {
	temp, err := template.ParseFiles(fileName)
	if err != nil {
		log.Fatal(err)
	}
	temp.Execute(w, info)
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
