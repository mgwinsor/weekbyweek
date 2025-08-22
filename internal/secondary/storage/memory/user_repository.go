package memory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

type userRepository struct {
	users map[uuid.UUID]*user.User
}

func NewUserRepository() user.UserRepository {
	return &userRepository{
		users: make(map[uuid.UUID]*user.User),
	}
}

func (r *userRepository) Save(ctx context.Context, user *user.User) error {
	r.users[user.ID()] = user
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	for _, user := range r.users {
		if user.Email() == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
