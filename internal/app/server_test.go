package app

import (
	"context"

	"github.com/mgwinsor/weekbyweek/internal/database"
)

type mockQuerier struct {
	database.Querier
	CreateUserFunc        func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByUsernameFunc func(ctx context.Context, username string) (database.User, error)
	GetUserByEmailFunc    func(ctx context.Context, email string) (database.User, error)
}

func (m *mockQuerier) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, arg)
	}
	panic("CreateUserFunc was called but not defined in mock")
}

func (m *mockQuerier) GetUserByUsername(ctx context.Context, username string) (database.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	panic("GetUserByUsernameFunc was called but not defined in mock")
}

func (m *mockQuerier) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	panic("GetUserByEmailFunc was called but not defined in mock")
}
