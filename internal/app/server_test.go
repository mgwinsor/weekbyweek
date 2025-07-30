package app

import (
	"context"
	"database/sql"

	"github.com/mgwinsor/weekbyweek/internal/database"
)

type mockDB struct {
	createUserResult           database.User
	createUserError            error
	deleteUserError            error
	getUserByUsernameResult    database.User
	getUserByUsernameError     error
	getUserByEmailResult       database.User
	getUserByEmailError        error
	updateUserDateofBirthError error
	updateUserEmailError       error
	updateUserPasswordError    error
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return m.createUserResult, m.createUserError
}

func (m *mockDB) DeleteUser(ctx context.Context, id interface{}) error {
	return m.deleteUserError
}

func (m *mockDB) GetUserByUsername(ctx context.Context, username string) (database.User, error) {
	return m.getUserByUsernameResult, m.getUserByUsernameError
}

func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return m.getUserByEmailResult, m.getUserByEmailError
}

func (m *mockDB) UpdateUserDateOfBirth(ctx context.Context, arg database.UpdateUserDateOfBirthParams) error {
	return m.updateUserDateofBirthError
}

func (m *mockDB) UpdateUserEmail(ctx context.Context, arg database.UpdateUserEmailParams) error {
	return m.updateUserEmailError
}

func (m *mockDB) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return m.updateUserPasswordError
}

func newDefaultMockDB() *mockDB {
	return &mockDB{
		createUserError:        nil,
		getUserByUsernameError: sql.ErrNoRows,
		getUserByEmailError:    sql.ErrNoRows,
	}
}

type mockAuthAdapter struct {
	hash string
	err  error
}

func (m *mockAuthAdapter) HashPassword(password string) (string, error) {
	return m.hash, m.err
}
