package user

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errHasherFailed = errors.New("hashing error")

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

func TestNewUser(t *testing.T) {
	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	futureDOB := time.Now().Add(24 * time.Hour)
	underageDOB := time.Now().Add(-10 * 365 * 24 * time.Hour)
	validPassword := "password"
	hashedPassword := "hashed-password"

	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		dateOfBirth time.Time
		mockSetup   func(m *MockPasswordHasher)
		expectedErr error
	}{
		{
			name:        "valid user",
			email:       "john@example.com",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: validDOB,
			mockSetup: func(m *MockPasswordHasher) {
				m.On("Hash", validPassword).
					Return(hashedPassword, nil).
					Once()
			},
			expectedErr: nil,
		},
		{
			name:        "empty email",
			email:       "",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: validDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrEmptyEmail,
		},
		{
			name:        "invalid email",
			email:       "invalid",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: validDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrInvalidEmailFormat,
		},
		{
			name:        "empty username",
			email:       "john@example.com",
			username:    "",
			password:    validPassword,
			dateOfBirth: validDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrEmptyUsername,
		},
		{
			name:        "future date of birth",
			email:       "john@example.com",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: futureDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrMininumAge,
		},
		{
			name:        "underage",
			email:       "john@example.com",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: underageDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrMininumAge,
		},
		{
			name:        "password too short",
			email:       "john@example.com",
			username:    "johndoe",
			password:    "pass",
			dateOfBirth: validDOB,
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrPasswordMinimumLenth,
		},
		{
			name:        "hashing error",
			email:       "john@example.com",
			username:    "johndoe",
			password:    validPassword,
			dateOfBirth: validDOB,
			mockSetup: func(m *MockPasswordHasher) {
				m.On("Hash", validPassword).
					Return("", errHasherFailed).
					Once()
			},
			expectedErr: errHasherFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHasher := new(MockPasswordHasher)
			tt.mockSetup(mockHasher)

			user, err := NewUser(tt.email, tt.username, tt.password, tt.dateOfBirth, mockHasher)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr, "expected an error but got none")
				assert.Nil(t, user, "user should be nil when error is returned")
			} else {
				require.NoError(t, err, "expected no error but got one")
				require.NotNil(t, user, "user should not be nil on success")

				assert.NotEqual(t, uuid.Nil, user.ID(), "expected a valid UUID, but it was nil")
				assert.Equal(t, tt.email, user.Email(), "email does not match expected")
				assert.Equal(t, tt.username, user.Username(), "username does not match expected")
				assert.NotEmpty(t, user.PasswordHash(), "password hash should be set on successful creation")
				assert.NotEqual(t, tt.password, user.PasswordHash(), "password hash should not be the same as the raw password")
				assert.Equal(t, tt.dateOfBirth, user.DateOfBirth(), "date of birth does not match expected")
			}
			mockHasher.AssertExpectations(t)
		})
	}
}

func TestNewUser_UniqueIDs(t *testing.T) {
	mockHasher := new(MockPasswordHasher)
	func(m *MockPasswordHasher) {
		m.On("Hash", "password").
			Return("hashed-password", nil).
			Twice()
	}(mockHasher)

	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	user1, _ := NewUser("john1@example.com", "John 1", "password", validDOB, mockHasher)
	user2, _ := NewUser("john2@example.com", "John 2", "password", validDOB, mockHasher)

	if user1.ID() == user2.ID() {
		t.Error("Expected users to have different IDs")
	}
}
