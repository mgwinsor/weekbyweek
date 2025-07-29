package auth

import (
	"errors"
	"testing"
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
