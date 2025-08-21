package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id          uuid.UUID
	email       string
	userName    string
	dateOfBirth time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewUser(email, userName string, dateOfBirth time.Time) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validateUserName(userName); err != nil {
		return nil, err
	}

	if err := validateDateOfBirth(dateOfBirth); err != nil {
		return nil, err
	}

	return &User{
		id:          uuid.New(),
		email:       email,
		userName:    userName,
		dateOfBirth: dateOfBirth,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}, nil
}

func (u *User) ID() uuid.UUID          { return u.id }
func (u *User) Email() string          { return u.email }
func (u *User) UserName() string       { return u.userName }
func (u *User) DateOfBirth() time.Time { return u.dateOfBirth }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}

	return nil
}

func validateUserName(userName string) error {
	if userName == "" {
		return errors.New("username cannot be empty")
	}

	return nil
}

func validateDateOfBirth(dob time.Time) error {
	if dob.IsZero() {
		return errors.New("date of birth is required")
	}

	age := time.Since(dob).Hours() / 24 / 365.25
	if age < 13 {
		return errors.New("user must be at least 13 years old")
	}

	if dob.After(time.Now()) {
		return errors.New("date of birth cannot be in the future")
	}

	return nil
}
