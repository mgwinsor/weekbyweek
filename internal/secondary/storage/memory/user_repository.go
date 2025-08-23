package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

type inMemoryUserRepository struct {
	users map[uuid.UUID]*user.User
	mu    sync.RWMutex
}

func NewUserRepository() user.UserRepository {
	return &inMemoryUserRepository{
		users: make(map[uuid.UUID]*user.User),
		mu:    sync.RWMutex{},
	}
}

func (r *inMemoryUserRepository) Save(ctx context.Context, user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID()] = user
	return nil
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id]
	if !exists {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (r *inMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, user.ErrUserNotFound
}
