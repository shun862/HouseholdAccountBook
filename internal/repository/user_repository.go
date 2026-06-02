package repository

import (
	"database/sql"
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
