package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

type mockUserRepository struct {
	saveError         error
	findByIDResult    *user.User
	findByIDError     error
	findByEmailResult *user.User
	findByEmailError  error
}

func newDefaultMockUserRepo() *mockUserRepository {
	return &mockUserRepository{
		saveError:         nil,
		findByIDResult:    nil,
		findByIDError:     user.ErrUserNotFound,
		findByEmailResult: nil,
		findByEmailError:  user.ErrUserNotFound,
	}
}

func (m *mockUserRepository) Save(ctx context.Context, u *user.User) error {
	return m.saveError
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return m.findByIDResult, m.findByIDError
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.findByEmailResult, m.findByEmailError
}

func TestCreateUser(t *testing.T) {
	dob := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	createUserRequest := CreateUserRequest{
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: dob,
	}

	tests := []struct {
		name     string
		req      CreateUserRequest
		mockRepo *mockUserRepository
		wantErr  bool
	}{
		{
			name:     "valid user",
			req:      createUserRequest,
			mockRepo: newDefaultMockUserRepo(),
			wantErr:  false,
		},
		{
			name: "duplicate email",
			req:  createUserRequest,
			mockRepo: &mockUserRepository{
				saveError:         nil,
				findByIDResult:    nil,
				findByIDError:     user.ErrUserNotFound,
				findByEmailResult: &user.User{},
				findByEmailError:  nil,
			},
			wantErr: true,
		},
		{
			name: "repository error during email lookup",
			req:  createUserRequest,
			mockRepo: &mockUserRepository{
				saveError:         nil,
				findByIDResult:    nil,
				findByIDError:     user.ErrUserNotFound,
				findByEmailResult: nil,
				findByEmailError:  errors.New("unknown error when finding email"),
			},
			wantErr: true,
		},
		{
			name: "repository error during save",
			req:  createUserRequest,
			mockRepo: &mockUserRepository{
				saveError:         errors.New("unknown error when saving user"),
				findByIDResult:    nil,
				findByIDError:     user.ErrUserNotFound,
				findByEmailResult: nil,
				findByEmailError:  user.ErrUserNotFound,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mockRepo)

			_, err := service.CreateUser(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
