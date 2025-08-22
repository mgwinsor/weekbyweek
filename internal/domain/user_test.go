package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	futureDOB := time.Now().Add(24 * time.Hour)
	underageDOB := time.Now().Add(-10 * 365 * 24 * time.Hour)

	tests := []struct {
		name        string
		email       string
		username    string
		dateOfBirth time.Time
		wantErr     bool
	}{
		{"valid user", "john@example.com", "John Doe", validDOB, false},
		{"invalid email", "invalid", "John", validDOB, true},
		{"empty username", "john@example.com", "", validDOB, true},
		{"future DOB", "john@example.com", "John", futureDOB, true},
		{"underage", "john@example.com", "John", underageDOB, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.username, tt.dateOfBirth)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {

				if user.ID() == uuid.Nil {
					t.Errorf("Expected valid UUID, got nil UUID")
				}
				if user.Email() != tt.email {
					t.Errorf("Expected email %s, got %s", tt.email, user.Email())
				}
				if user.Username() != tt.username {
					t.Errorf("Expected username %s, got %s", tt.username, user.Username())
				}
				if user.DateOfBirth() != tt.dateOfBirth {
					t.Errorf("Expected date of birth %v, got %v", validDOB, user.DateOfBirth())
				}
			}
		})
	}
}

func TestNewUser_UniqueIDs(t *testing.T) {
	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	user1, _ := NewUser("john1@example.com", "John 1", validDOB)
	user2, _ := NewUser("john2@example.com", "John 2", validDOB)

	if user1.ID() == user2.ID() {
		t.Error("Expected users to have different IDs")
	}
}
