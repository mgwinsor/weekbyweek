package app

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterPost(t *testing.T) {

	templates := template.Must(template.ParseGlob("../../web/templates/*.html"))

	tests := []struct {
		name                string
		mockAuth            *mockAuthAdapter
		mockDB              *mockDB
		formInput           string
		wantCode            int
		wantBodyContains    string
		notWantBodyContains string
	}{
		{
			name: "successfully creates user",
			mockAuth: &mockAuthAdapter{
				hash: "hashed-password",
				err:  nil,
			},
			mockDB:              newDefaultMockDB(),
			formInput:           "username=newuser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			wantCode:            http.StatusCreated,
			wantBodyContains:    "",
			notWantBodyContains: "Create Your Account",
		},
		{
			name: "fails to create user with duplicate username",
			mockDB: func() *mockDB {
				m := newDefaultMockDB()
				m.getUserByUsernameError = nil
				return m
			}(),
			formInput:        "username=existinguser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			wantCode:         http.StatusOK,
			wantBodyContains: ErrorUsernameExists.Error(),
		},
		{
			name: "fails to create user with duplicate email",
			mockDB: func() *mockDB {
				m := newDefaultMockDB()
				m.getUserByEmailError = nil
				return m
			}(),
			formInput:        "username=newuser&email=existing@example.com&password=A$$w0rd12345&dob=1992-11-21",
			wantCode:         http.StatusOK,
			wantBodyContains: ErrorEmailExists.Error(),
		},
		{
			name:             "fails to validate weak password",
			mockDB:           newDefaultMockDB(),
			formInput:        "username=newuser&email=existing@example.com&password=ab&dob=1992-11-21",
			wantCode:         http.StatusOK,
			wantBodyContains: ErrorPasswordTooShort.Error(),
		},
		{
			name: "fails to create user in database",
			mockAuth: &mockAuthAdapter{
				hash: "hashed-password",
				err:  nil,
			},
			mockDB: func() *mockDB {
				m := newDefaultMockDB()
				m.createUserError = errors.New("Internal database error")
				return m
			}(),
			formInput:        "username=newuser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			wantCode:         http.StatusInternalServerError,
			wantBodyContains: "Internal Server Error",
		},
		{
			name: "fails to hash password",
			mockAuth: &mockAuthAdapter{
				hash: "",
				err:  errors.New("bcrypt error"),
			},
			mockDB:           newDefaultMockDB(),
			formInput:        "username=newuser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			wantCode:         http.StatusInternalServerError,
			wantBodyContains: "Could not hash password",
		},
		{
			name:      "fails to get data from database with username",
			formInput: "username=newuser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			mockDB: func() *mockDB {
				m := newDefaultMockDB()
				m.getUserByUsernameError = errors.New("Internal database error")
				return m
			}(),
			wantCode:         http.StatusInternalServerError,
			wantBodyContains: "Internal Server Error",
		},
		{
			name:      "fails to get data from database with email",
			formInput: "username=newuser&email=new@example.com&password=A$$w0rd12345&dob=1992-11-21",
			mockDB: func() *mockDB {
				m := newDefaultMockDB()
				m.getUserByEmailError = errors.New("Internal database error")
				return m
			}(),
			wantCode:         http.StatusInternalServerError,
			wantBodyContains: "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := NewServer(tt.mockDB, tt.mockAuth, templates)
			form := strings.NewReader(tt.formInput)
			req := httptest.NewRequest("POST", "/register", form)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			svr.registerPost(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status code %d; got %d", tt.wantCode, rr.Code)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tt.wantBodyContains) {
				t.Errorf("expected body to contain %s; it did not", tt.wantBodyContains)
			}

			if tt.notWantBodyContains != "" && strings.Contains(body, tt.notWantBodyContains) {
				t.Errorf("expected body NOT to contain %s; it did", tt.notWantBodyContains)
			}
		})
	}
}
