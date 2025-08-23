package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
