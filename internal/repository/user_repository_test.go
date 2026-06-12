package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func setupTestUserTable(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_name TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL
        );
    `)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestInsertIntoUser(t *testing.T) {
	db := setupTestUserTable(t)
	defer db.Close()

	repo := NewUserRepository(db)
	err := repo.Insert("test", "password")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	rows, err := db.Query("SELECT user_name, password FROM users WHERE user_name = ?", "test")
	if err != nil {
		t.Fatalf("select query failed: %v", err)
	}

	defer rows.Close()

	count := 0
	var n string
	var p string
	for rows.Next() {
		rows.Scan(&n, &p)
		count++
	}

	if n != "test" || p != "password" {
		t.Errorf("insert value is wrong. username:%s password:%s", n, p)
	}

	if count != 1 {
		t.Errorf("expected 1 user. count:%d", count)
	}
}

func TestCheckSameUser(t *testing.T) {
	db := setupTestUserTable(t)
	defer db.Close()

	repo := NewUserRepository(db)
	repo.Insert("test", "password")
	hasUser, err := repo.CheckSameUser("test")
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}

	if !hasUser {
		t.Errorf("expected user to exist")
	}
}

func TestCheckSameUser_NotFound(t *testing.T) {
	db := setupTestUserTable(t)
	defer db.Close()

	repo := NewUserRepository(db)
	repo.Insert("test", "password")
	hasUser, err := repo.CheckSameUser("tes")
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}

	if hasUser {
		t.Errorf("expected user not to exist")
	}
}

func TestFindUser(t *testing.T) {
	db := setupTestUserTable(t)
	repo := NewUserRepository(db)
	password, _ := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)
	repo.Insert("test", string(password))
	user, err := repo.FindUser("test", "password")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}

	if user == nil {
		t.Fatalf("user not find.")
	}

	if user.UserName != "test" {
		t.Errorf("username is wrong.")
	}
}

func TestFindUser_WrongPassword(t *testing.T) {
	db := setupTestUserTable(t)
	repo := NewUserRepository(db)
	password, _ := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)
	repo.Insert("test", string(password))
	user, err := repo.FindUser("test", "pass")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}

	if user != nil {
		t.Errorf("expected nil user. username:%s password:%s",
			user.UserName, user.Password)
	}
}

func TestFindUser_NotFound(t *testing.T) {
	db := setupTestUserTable(t)
	repo := NewUserRepository(db)
	password, _ := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)
	repo.Insert("test", string(password))
	user, err := repo.FindUser("tes", "password")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}

	if user != nil {
		t.Errorf("expected nil user. username:%s password:%s",
			user.UserName, user.Password)
	}
}
