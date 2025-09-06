package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrEmptyEmail = errors.New("email cannot be empty")
var ErrInvalidEmailFormat = errors.New("incorrect email format")
var ErrEmptyUsername = errors.New("username cannot be empty")
var ErrPasswordMinimumLenth = errors.New("password must be at least 8 characters long")
var ErrMininumAge = errors.New("user must be at least 13 years old")

type User struct {
	id           uuid.UUID
	email        string
	username     string
	passwordHash string
	dateOfBirth  time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(email, username, password string, dateOfBirth time.Time, hasher PasswordHasher) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if err := validateDateOfBirth(dateOfBirth); err != nil {
		return nil, err
	}

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           uuid.New(),
		email:        email,
		username:     username,
		passwordHash: hashedPassword,
		dateOfBirth:  dateOfBirth,
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
		return ErrEmptyEmail
	}

	if !strings.Contains(email, "@") {
		return ErrInvalidEmailFormat
	}

	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordMinimumLenth
	}
	return nil
}

func validateDateOfBirth(dob time.Time) error {
	age := time.Since(dob).Hours() / 24 / 365.25
	if age < 13 {
		return ErrMininumAge
	}

	return nil
}
