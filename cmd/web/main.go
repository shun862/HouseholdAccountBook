package main

import (
	"fmt"
	"household_account_book/internal/db"
	"household_account_book/internal/handler"
	"household_account_book/internal/repository"
	"log"
	"net/http"

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

	// 静的ファイル配信
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../../web/css"))))

	// ルーティング
	http.HandleFunc("/user_register", userHandler.HandleFunc)

	fmt.Println("server start :8080")
	// サーバー起動
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
