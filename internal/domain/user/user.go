package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmailRequired      = errors.New("email cannot be empty")
	ErrInvalidEmailFormat = errors.New("incorrect email format")
	ErrUsernameRequired   = errors.New("username cannot be empty")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
)

type NewUserParams struct {
	Email       string
	Username    string
	Password    string
	DateOfBirth time.Time
}

type User struct {
	id           uuid.UUID
	email        string
	username     string
	passwordHash string
	dateOfBirth  time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(params NewUserParams, hasher PasswordHasher) (*User, error) {
	if err := validateEmail(params.Email); err != nil {
		return nil, err
	}

	if err := validateUsername(params.Username); err != nil {
		return nil, err
	}

	if err := validatePassword(params.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := hasher.Hash(params.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           uuid.New(),
		email:        params.Email,
		username:     params.Username,
		passwordHash: hashedPassword,
		dateOfBirth:  params.DateOfBirth,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}, nil
}

func (u *User) ID() uuid.UUID          { return u.id }
func (u *User) Email() string          { return u.email }
func (u *User) Username() string       { return u.username }
func (u *User) PasswordHash() string   { return u.passwordHash }
func (u *User) DateOfBirth() time.Time { return u.dateOfBirth }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }

func validateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	if !strings.Contains(email, "@") {
		return ErrInvalidEmailFormat
	}

	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return ErrUsernameRequired
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
