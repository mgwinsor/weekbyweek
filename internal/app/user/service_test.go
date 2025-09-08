package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errRepositoryFailure = errors.New("error in data repository")

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}
	return u, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}
	return u, args.Error(1)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	dob := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	createUserRequest := CreateUserRequest{
		Email:       "john@example.com",
		Username:    "johndoe",
		Password:    "12345678",
		DateOfBirth: dob,
	}

	setupHasher := new(MockPasswordHasher)
	setupHasher.On("Hash", "password").Return("hashed-password", nil)
	existingUser, _ := user.NewUser(
		user.NewUserParams{
			Email:       "john@example.com",
			Username:    "existing-user",
			Password:    "password",
			DateOfBirth: dob,
		},
		setupHasher,
	)

	tests := []struct {
		name        string
		req         CreateUserRequest
		mockSetup   func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher)
		expectedErr error
	}{
		{
			name: "successfully create user",
			req:  createUserRequest,
			mockSetup: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				mockHasher.On("Hash", createUserRequest.Password).
					Return("hashed-password", nil).Once()
				mockRepo.On("FindByEmail", mock.Anything, createUserRequest.Email).
					Return(nil, user.ErrUserNotFound).Once()
				mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "error on duplicate email",
			req:  createUserRequest,
			mockSetup: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				mockRepo.On("FindByEmail", mock.Anything, createUserRequest.Email).
					Return(existingUser, nil).Once()
			},
			expectedErr: ErrEmailExists,
		},
		{
			name: "repository error during email lookup",
			req:  createUserRequest,
			mockSetup: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				mockRepo.On("FindByEmail", mock.Anything, createUserRequest.Email).
					Return(nil, errRepositoryFailure).Once()
			},
			expectedErr: errRepositoryFailure,
		},
		{
			name: "repository error during save",
			req:  createUserRequest,
			mockSetup: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				mockHasher.On("Hash", createUserRequest.Password).
					Return("hashed-password", nil).Once()
				mockRepo.On("FindByEmail", mock.Anything, createUserRequest.Email).
					Return(nil, user.ErrUserNotFound).Once()
				mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).
					Return(errRepositoryFailure)
			},
			expectedErr: errRepositoryFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockHasher := new(MockPasswordHasher)
			tt.mockSetup(mockRepo, mockHasher)

			userService := NewUserService(mockRepo, mockHasher)

			resp, err := userService.CreateUser(context.Background(), tt.req)

			if tt.expectedErr != nil {
				require.Error(t, err)

				if errors.Is(tt.expectedErr, ErrEmailExists) {
					assert.ErrorIs(t, err, tt.expectedErr)
				} else if errors.Is(tt.expectedErr, errRepositoryFailure) {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, resp, "response should be nil when error is returned")
			} else {
				require.NoError(t, err, "CreateUser failed unexpectedly")
				require.NotNil(t, resp, "response should not be nil on success")
			}
		})
	}
}
