package app

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/database"
)

func TestRegisterPost(t *testing.T) {

	templates := template.Must(template.ParseGlob("../../web/templates/*.html"))

	successMockDB := &mockQuerier{
		CreateUserFunc: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
			id, _ := uuid.Parse("1216461d-0f93-4059-98c3-bf922acb752b")
			return database.User{ID: id, Username: arg.Username, Email: arg.Email}, nil
		},
		GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{}, sql.ErrNoRows
		},
		GetUserByUsernameFunc: func(ctx context.Context, username string) (database.User, error) {
			return database.User{}, sql.ErrNoRows
		},
	}

	failUsernameMockDB := &mockQuerier{
		GetUserByUsernameFunc: func(ctx context.Context, username string) (database.User, error) {
			return database.User{Username: "existinguser"}, nil
		},
	}

	failEmailMockDB := &mockQuerier{
		GetUserByUsernameFunc: func(ctx context.Context, username string) (database.User, error) {
			return database.User{}, sql.ErrNoRows
		},
		GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{Email: "existing@example.com"}, nil
		},
	}

	tests := []struct {
		name     string
		mockDB   *mockQuerier
		formData string
		wantCode int
	}{
		{
			name:     "successfully creates user",
			mockDB:   successMockDB,
			formData: "username=newuser&email=new@example.com&password=a$$w0rd123&dob=1992-11-21",
			wantCode: http.StatusCreated,
		},
		{
			name:     "fails to create user with duplicate username",
			mockDB:   failUsernameMockDB,
			formData: "username=existinguser&email=new@example.com&password=a$$w0rd123&dob=1992-11-21",
			wantCode: http.StatusConflict,
		},
		{
			name:     "fails to create user with duplicate email",
			mockDB:   failEmailMockDB,
			formData: "username=newuser&email=existing@example.com&password=a$$w0rd123&dob=1992-11-21",
			wantCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := NewServer(tt.mockDB, templates)
			form := strings.NewReader(tt.formData)
			req := httptest.NewRequest("POST", "/register", form)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			svr.registerPost(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status code %d; got %d", tt.wantCode, rr.Code)
			}
		})
	}
}
