package memory

import (
	"context"
	"testing"
	"time"

	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

func TestUserRepository(t *testing.T) {
	repo := NewUserRepository()
	dob := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	user, _ := user.NewUser("john@example.com", "John Doe", dob)

	err := repo.Save(context.Background(), user)
	if err != nil {
		t.Fatalf("Expected no error saving user, got %v", err)
	}

	foundUser, err := repo.FindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("Expected to find user, got %v", err)
	}

	if foundUser.Email() != user.Email() {
		t.Errorf("Expected email %s, got %s", user.Email(), foundUser.Email())
	}
}
