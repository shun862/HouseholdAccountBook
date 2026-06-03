package repository

import (
	"database/sql"
	"household_account_book/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) Insert(username string, password string) error {
	query := `
	INSERT INTO users (user_name, password)
	VALUES(?, ?)
	`

	_, err := repo.DB.Exec(query, username, password)
	return err
}

func (repo *UserRepository) CheckSameUser(username string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM users WHERE user_name = ?
	`
	var count int
	err := repo.DB.QueryRow(query, username).Scan(&count)

	hasSameUser := count > 0
	return hasSameUser, err
}

func (repo *UserRepository) FindUser(username string, password string) (*model.User, error) {
	query := `
	SELECT id, user_name, password FROM users WHERE user_name = ?
	`
	var user model.User
	err := repo.DB.QueryRow(query, username).Scan(&user.Id, &user.UserName, &user.Password)
	if err != nil {
		return nil, err
	}

	// パスワード照合
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// パスワード不一致
		return nil, nil
	}

	return &user, nil
}
