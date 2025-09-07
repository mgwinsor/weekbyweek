package memory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHasher struct{}

func (f *fakeHasher) Hash(password string) (string, error)          { return "hashed-" + password, nil }
func (f *fakeHasher) Compare(hashedPassword, password string) error { return nil }

func TestUserRepository(t *testing.T) {
	validUser, err := user.NewUser(
		user.NewUserParams{
			Email:       "john@example.com",
			Username:    "johndoe",
			Password:    "password",
			DateOfBirth: time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC),
		},
		&fakeHasher{},
	)
	require.NoError(t, err, "failed to create test user")

	t.Run("save and find user by ID and email", func(t *testing.T) {
		repo := NewUserRepository()

		err := repo.Save(context.Background(), validUser)
		require.NoError(t, err)

		foundByID, err := repo.FindByID(context.Background(), validUser.ID())
		require.NoError(t, err)

		foundByEmail, err := repo.FindByEmail(context.Background(), validUser.Email())
		require.NoError(t, err)

		assert.Equal(t, validUser, foundByID, "user found by ID should match the saved user")
		assert.Equal(t, validUser, foundByEmail, "user found by email should match the saved user")
	})

	t.Run("return error for non-existent ID", func(t *testing.T) {
		repo := NewUserRepository()

		foundUser, err := repo.FindByID(context.Background(), uuid.New())

		require.Error(t, err, "expected an error for non-existent user")
		assert.ErrorIs(t, err, user.ErrUserNotFound, "error should be ErrUserNotFound")
		assert.Nil(t, foundUser, "found user should be nil on error")
	})

	t.Run("return error for non-existent email", func(t *testing.T) {
		repo := NewUserRepository()

		foundUser, err := repo.FindByEmail(context.Background(), "notfound@example.com")

		require.Error(t, err, "expected an error for non-existent user")
		assert.ErrorIs(t, err, user.ErrUserNotFound, "error should be ErrUserNotFound")
		assert.Nil(t, foundUser, "found user should be nil on error")
	})
}
