package handler

import (
	"context"
	"household_account_book/internal/consts"
	"household_account_book/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type MockRegisterUserRepositoryWithSameUser struct {
}

func (m *MockRegisterUserRepositoryWithSameUser) Insert(username, password string) error {
	return nil
}

func (m *MockRegisterUserRepositoryWithSameUser) CheckSameUser(username string) (bool, error) {
	return true, nil
}

func (m *MockRegisterUserRepositoryWithSameUser) FindUser(username, password string) (*model.User, error) {
	return &model.User{
		Id:       1,
		UserName: username,
	}, nil
}

type MockRegisterUserRepository struct {
}

func (m *MockRegisterUserRepository) Insert(username, password string) error {
	return nil
}

func (m *MockRegisterUserRepository) CheckSameUser(username string) (bool, error) {
	return false, nil
}

func (m *MockRegisterUserRepository) FindUser(username, password string) (*model.User, error) {
	return &model.User{
		Id:       1,
		UserName: username,
	}, nil
}

func TestRegisterHandleFunc_Get(t *testing.T) {
	handler := NewUserHandler(&MockRegisterUserRepository{})
	req := httptest.NewRequest(
		http.MethodGet,
		consts.UserRegisterUrl,
		nil,
	)
	rr := httptest.NewRecorder()
	handler.RegisterHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d", rr.Code)
	}
}

func TestRegisterHandleFunc_Post(t *testing.T) {
	handler := NewUserHandler(&MockRegisterUserRepository{})

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "password123")

	req := httptest.NewRequest(
		http.MethodPost,
		consts.UserRegisterUrl,
		strings.NewReader(form.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rr := httptest.NewRecorder()
	handler.RegisterHandleFunc(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("status = %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != consts.LoginUrl {
		t.Errorf(
			"location = %s", location,
		)
	}
}

func TestRegisterHandleFunc_Post_ValidationError(t *testing.T) {
	handler := NewUserHandler(&MockRegisterUserRepository{})

	form := url.Values{}
	form.Set("username", "")
	form.Set("password", "password123")

	req := httptest.NewRequest(
		http.MethodPost,
		consts.UserRegisterUrl,
		strings.NewReader(form.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rr := httptest.NewRecorder()
	handler.RegisterHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d", rr.Code)
	}
}

func TestRegisterHandleFunc_Post_SameUserError(t *testing.T) {
	handler := NewUserHandler(&MockRegisterUserRepositoryWithSameUser{})

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "password123")

	req := httptest.NewRequest(
		http.MethodPost,
		consts.UserRegisterUrl,
		strings.NewReader(form.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rr := httptest.NewRecorder()
	handler.RegisterHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "すでに登録済みのユーザーです") {
		t.Errorf("expected error message")
	}
}

type MockLoginUserRepository struct {
}

func (m *MockLoginUserRepository) Insert(username, password string) error {
	return nil
}

func (m *MockLoginUserRepository) CheckSameUser(username string) (bool, error) {
	return false, nil
}

func (m *MockLoginUserRepository) FindUser(username, password string) (*model.User, error) {
	return &model.User{
		Id:       1,
		UserName: username,
	}, nil
}

type MockLoginUserRepositoryNotFoundUser struct {
}

func (m *MockLoginUserRepositoryNotFoundUser) Insert(username, password string) error {
	return nil
}

func (m *MockLoginUserRepositoryNotFoundUser) CheckSameUser(username string) (bool, error) {
	return false, nil
}

func (m *MockLoginUserRepositoryNotFoundUser) FindUser(username, password string) (*model.User, error) {
	return nil, nil
}

func TestLoginHandleFunc_Get(t *testing.T) {
	handler := NewUserHandler(&MockLoginUserRepository{})

	req := httptest.NewRequest(http.MethodGet, consts.LoginUrl, nil)
	rr := httptest.NewRecorder()

	handler.LoginHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d", rr.Code)
	}
}

func TestLoginHandleFunc_Post(t *testing.T) {
	handler := NewUserHandler(&MockLoginUserRepository{})

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, consts.LoginUrl, strings.NewReader(form.Encode()))
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	rr := httptest.NewRecorder()

	handler.LoginHandleFunc(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("status = %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != consts.ExpenseListUrl {
		t.Errorf(
			"location = %s", location,
		)
	}

	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == consts.TokenName {
			found = true
			if cookie.Value == "" {
				t.Error("token cookie is empty")
			}
		}
	}

	if !found {
		t.Error("token cookie not found")
	}
}

func TestLoginHandleFunc_Post_ValidationError(t *testing.T) {
	handler := NewUserHandler(&MockLoginUserRepository{})

	form := url.Values{}
	form.Set("username", "")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, consts.LoginUrl, strings.NewReader(form.Encode()))
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	rr := httptest.NewRecorder()

	handler.LoginHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d", rr.Code)
	}
}

func TestLoginHandleFunc_Post_NotFoundUser(t *testing.T) {
	handler := NewUserHandler(&MockLoginUserRepositoryNotFoundUser{})

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, consts.LoginUrl, strings.NewReader(form.Encode()))
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	rr := httptest.NewRecorder()

	handler.LoginHandleFunc(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "未登録のユーザーです") {
		t.Errorf("expected error message")
	}
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		expectError   bool
		usernameError string
		passwordError string
	}{
		{
			name:        "正常",
			username:    "testuser",
			password:    "password",
			expectError: false,
		},
		{
			name:          "ユーザー名未入力",
			username:      "",
			password:      "password",
			expectError:   true,
			usernameError: "ユーザー名を入力してください",
		},
		{
			name:          "ユーザー名文字数不足",
			username:      "abc",
			password:      "password",
			expectError:   true,
			usernameError: "ユーザー名は4～10文字で入力してください",
		},
		{
			name:          "パスワード未入力",
			username:      "tester",
			password:      "",
			expectError:   true,
			passwordError: "パスワードを入力してください",
		},
		{
			name:          "パスワード文字数不足",
			username:      "tester",
			password:      "1234567",
			expectError:   true,
			passwordError: "パスワードは8～12文字で入力してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := UserViewInfo{}

			userValidation(tt.username, tt.password, &info)

			if info.IsError != tt.expectError {
				t.Errorf("want %v got %v", tt.expectError, info.IsError)
			}

			if info.UsernameError != tt.usernameError {
				t.Errorf("want %s got %s",
					tt.usernameError,
					info.UsernameError)
			}

			if info.PasswordError != tt.passwordError {
				t.Errorf("want %s got %s",
					tt.passwordError,
					info.PasswordError)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		consts.UserIDKey,
		123,
	)

	id := GetUserID(ctx)
	if id != 123 {
		t.Errorf("want 123 got %d", id)
	}
}

func TestGetUserID_NotFound(t *testing.T) {
	id := GetUserID(context.Background())
	if id != 0 {
		t.Errorf("want 0 got %d", id)
	}
}

func TestLogoutHandleFunc(t *testing.T) {
	handler := &UserHandler{}

	req := httptest.NewRequest(
		http.MethodGet,
		consts.LoginUrl,
		nil,
	)
	rec := httptest.NewRecorder()
	handler.LogoutHandleFunc(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusSeeOther {
		t.Errorf("want %d got %d", http.StatusSeeOther, res.StatusCode)
	}

	cookies := res.Cookies()
	if len(cookies) == 0 {
		t.Fatal("cookie not found")
	}

	if cookies[0].Name != consts.TokenName {
		t.Errorf("unexpected cookie")
	}

	if cookies[0].Value != "" {
		t.Errorf("cookie should be empty")
	}
}
