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

const (
	validPassword  = "12345678"
	hashedPassword = "hashed-password"
)

var (
	errHasherFailed    = errors.New("hashing error")
	validDOB           = time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	validNewUserParams = NewUserParams{
		Email:       "john@example.com",
		Username:    "johndoe",
		Password:    validPassword,
		DateOfBirth: validDOB,
	}
)

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

func withParams(modifier func(p *NewUserParams)) NewUserParams {
	params := validNewUserParams
	modifier(&params)
	return params
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		params      NewUserParams
		mockSetup   func(m *MockPasswordHasher)
		expectedErr error
	}{
		{
			name:   "valid user",
			params: validNewUserParams,
			mockSetup: func(m *MockPasswordHasher) {
				m.On("Hash", validPassword).
					Return(hashedPassword, nil).
					Once()
			},
			expectedErr: nil,
		},
		{
			name:        "empty email",
			params:      withParams(func(p *NewUserParams) { p.Email = "" }),
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrEmailRequired,
		},
		{
			name:        "invalid email",
			params:      withParams(func(p *NewUserParams) { p.Email = "invalid" }),
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrInvalidEmailFormat,
		},
		{
			name:        "empty username",
			params:      withParams(func(p *NewUserParams) { p.Username = "" }),
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrUsernameRequired,
		},
		{
			name:   "password is minimum length",
			params: validNewUserParams,
			mockSetup: func(m *MockPasswordHasher) {
				m.On("Hash", validPassword).
					Return(hashedPassword, nil).
					Once()
			},
			expectedErr: nil,
		},
		{
			name:        "password too short",
			params:      withParams(func(p *NewUserParams) { p.Password = "1234567" }),
			mockSetup:   func(m *MockPasswordHasher) {},
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:   "hashing error",
			params: validNewUserParams,
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

			user, err := NewUser(tt.params, mockHasher)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr, "expected an error but got none")
				assert.Nil(t, user, "user should be nil when error is returned")
			} else {
				require.NoError(t, err, "expected no error but got one")
				require.NotNil(t, user, "user should not be nil on success")

				assert.NotEqual(t, uuid.Nil, user.ID(), "expected a valid UUID, but it was nil")
				assert.Equal(t, tt.params.Email, user.Email(), "email does not match expected")
				assert.Equal(t, tt.params.Username, user.Username(), "username does not match expected")
				assert.NotEmpty(t, user.PasswordHash(), "password hash should be set on successful creation")
				assert.NotEqual(t, tt.params.Password, user.PasswordHash(), "password hash should not be the same as the raw password")
				assert.Equal(t, tt.params.DateOfBirth, user.DateOfBirth(), "date of birth does not match expected")
			}
			mockHasher.AssertExpectations(t)
		})
	}
}

func TestNewUser_UniqueIDs(t *testing.T) {
	mockHasher := new(MockPasswordHasher)
	mockHasher.On("Hash", validPassword).Return(hashedPassword, nil).Twice()

	params1 := validNewUserParams
	params1.Email = "john1@example.com"

	params2 := validNewUserParams
	params2.Email = "john2@example.com"

	user1, err1 := NewUser(params1, mockHasher)
	require.NoError(t, err1)

	user2, err2 := NewUser(params2, mockHasher)
	require.NoError(t, err2)

	assert.NotEqual(t, user1.ID(), user2.ID(), "expected users to have different IDs")
}
