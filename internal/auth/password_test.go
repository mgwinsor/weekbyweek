package auth

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockHasher struct {
	hash []byte
	err  error
}

func (m *mockHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return m.hash, m.err
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		hasher    PasswordHasher
		wantHash  string
		wantError bool
	}{
		{
			name: "successful hash",
			hasher: &mockHasher{
				hash: []byte("hashedpassword"),
				err:  nil,
			},
			wantHash:  "hashedpassword",
			wantError: false,
		},
		{
			name: "hashing error",
			hasher: &mockHasher{
				hash: nil,
				err:  errors.New("bcrypt error"),
			},
			wantHash:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword("password123", tt.hasher)

			if hash != tt.wantHash {
				t.Errorf("want hash %q; got %q", tt.wantHash, hash)
			}

			if (err != nil) != tt.wantError {
				t.Errorf("want error: %v; got %v", tt.wantError, err)
			}
		})
	}
}

func TestHashPassword_Integration(t *testing.T) {
	password := "Password123!"
	hasher := BcryptHasher{}

	hash, err := HashPassword(password, hasher)
	if err != nil {
		t.Errorf("HashPassword() returned an unexpected error: %v", err)
	}
	if hash == "" {
		t.Errorf("HashPassword() returned an empty hash")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Fatalf("bcrypt.CompareHashAndPassword() failed: %v", err)
	}
}
