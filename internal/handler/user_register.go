package handler

import (
	"household_account_book/internal/repository"
	"log"
	"net/http"
	"text/template"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (handler *UserHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		temp, err := template.ParseFiles("../../web/templates/user_register.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, nil)
	}
}
