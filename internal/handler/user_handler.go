package handler

import (
	"context"
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

type UserViewInfo struct {
	Username      string
	UsernameError string
	PasswordError string
	IsError       bool
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (handler *UserHandler) RegisterHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		info := userRegisterValidation(username, password, handler.Repo)
		if info.IsError {
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
		showView(w, consts.UserRegisterFile, UserViewInfo{})
	}
}

func (handler *UserHandler) LoginHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		info := UserViewInfo{
			Username: username,
		}
		userValidation(username, password, &info)
		if info.IsError {
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
			info.IsError = true
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
			Path:     consts.CookieValidPath,
		})

		http.Redirect(w, r, consts.AddExpenseUrl, http.StatusSeeOther)
	} else {
		showView(w, consts.LoginFile, UserViewInfo{})
	}
}

func showView(w http.ResponseWriter, fileName string, info UserViewInfo) {
	temp, err := template.ParseFiles(fileName)
	if err != nil {
		log.Fatal(err)
	}
	temp.Execute(w, info)
}

func (handler *UserHandler) LogoutHandleFunc(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     consts.TokenName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     consts.CookieValidPath,
		HttpOnly: true,
	})
	http.Redirect(w, r, consts.LoginUrl, http.StatusSeeOther)
}

func userRegisterValidation(username string, password string, repo *repository.UserRepository) UserViewInfo {
	info := UserViewInfo{
		Username: username,
	}

	// 共通バリデーション
	userValidation(username, password, &info)

	if info.IsError {
		return info
	}

	return info
}

// 共通バリデーション
func userValidation(username string, password string, info *UserViewInfo) {
	// ユーザー名
	if username == "" {
		info.UsernameError = "ユーザー名を入力してください"
		info.IsError = true
	} else if len := len([]rune(username)); len < 4 || len > 10 {
		info.UsernameError = "ユーザー名は4～10文字で入力してください"
		info.IsError = true
	}

	// パスワード
	if password == "" {
		info.PasswordError = "パスワードを入力してください"
		info.IsError = true
	} else if len := len([]rune(password)); len < 8 || len > 12 {
		info.PasswordError = "パスワードは8～12文字で入力してください"
		info.IsError = true
	}
}

func GetUserID(ctx context.Context) int {
	if v := ctx.Value(consts.UserIDKey); v != nil {
		return v.(int)
	}
	return 0
}
