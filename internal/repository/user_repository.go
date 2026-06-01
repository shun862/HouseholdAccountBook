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

func (repo *UserRepository) Insert(userName string, password string) error {
	query := `
	INSERT INTO users (user_name, password)
	VALUES(?, ?), 
	`

	_, err := repo.DB.Exec(query, userName, password)
	return err
}
