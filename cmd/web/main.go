package main

import (
	"fmt"
	"household_account_book/internal/consts"
	"household_account_book/internal/db"
	"household_account_book/internal/handler"
	"household_account_book/internal/model"
	"household_account_book/internal/repository"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// DB取得
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	// 最後にDBを閉じる
	defer db.Close()

	// DB接続確認
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)
	expenseRepo := repository.NewExpenseRepository(db)
	addExpenseHandler := handler.NewAddExpenseHandler(expenseRepo)

	// 静的ファイル配信
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../../web/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../../web/js"))))
	// ルーティング
	http.HandleFunc(consts.UserRegisterUrl, userHandler.RegisterHandleFunc)
	http.HandleFunc(consts.LoginUrl, userHandler.LoginHandleFunc)
	http.HandleFunc(consts.LogoutUrl, userHandler.LogoutHandleFunc)
	http.Handle(consts.AddExpenseUrl, authMiddleware(http.HandlerFunc(addExpenseHandler.HandleFunc)))

	fmt.Println("server start :8080")
	// サーバー起動
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(consts.TokenName)
		// 未ログインの場合
		if err != nil {
			http.Redirect(w, r, consts.LoginUrl, http.StatusSeeOther)
			return
		}

		tokenStr := cookie.Value
		claims := &model.JwtCustomClaims{}
		// トークン情報の取得
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(consts.JwtKey), nil
		})

		// トークンが無効
		if err != nil || !token.Valid {
			http.Redirect(w, r, consts.LoginUrl, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
