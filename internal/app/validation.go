package app

import (
	"errors"
	"net/mail"
	"time"
	"unicode"
)

var (
	ErrorDateOfBirthInvalid = errors.New("Must provide a valid date of birth in format YYYY-MM-DD")
	ErrorDateOfBirthInFutre = errors.New("Date of birth cannot be in the future")
	ErrorDateOfBirthTooOld  = errors.New("Date of birth cannot be too far in the past")

	ErrorUsernameTooShort     = errors.New("Username must be between 3 and 21 characters")
	ErrorUsernameTooLong      = errors.New("Username must be between 3 and 21 characters")
	ErrorUsernameInvalidChars = errors.New("Username can only contain letters, numbers, and underscores")
	ErrorUsernameExists       = errors.New("Username already exists")

	ErrorEmailInvalid = errors.New("Email is invalid")
	ErrorEmailExists  = errors.New("Email already exists")

	ErrorPasswordTooShort           = errors.New("Password must be at least 12 characters long")
	ErrorPasswordMissingUppercase   = errors.New("Password must contain at least 1 uppercase letter")
	ErrorPasswordMissingLowercase   = errors.New("Password must contain at least 1 lowercase letter")
	ErrorPasswordMissingNumber      = errors.New("Password must contain at least 1 number")
	ErrorPasswordMissingSpecialChar = errors.New("Password must contain at least 1 special character")
)

func validateDateOfBirth(dateOfBirth string) (time.Time, error) {
	dob, err := time.Parse(time.DateOnly, dateOfBirth)
	if err != nil {
		return time.Time{}, ErrorDateOfBirthInvalid
	}

	if dob.After(time.Now()) {
		return time.Time{}, ErrorDateOfBirthInFutre
	}

	if dob.Before(time.Now().AddDate(-90, 0, 0)) {
		return time.Time{}, ErrorDateOfBirthTooOld
	}

	return dob, nil
}

func validateUsernameSyntax(username string) error {
	if len(username) < 3 {
		return ErrorUsernameTooShort
	}
	if len(username) > 21 {
		return ErrorUsernameTooLong
	}

	for _, c := range username {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_' {
			return ErrorUsernameInvalidChars
		}
	}

	return nil
}

func validateEmailSyntax(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrorEmailInvalid
	}
	return nil
}

func valdiatePassword(password string) error {
	if len(password) < 12 {
		return ErrorPasswordTooShort
	}

	var (
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSpecialChar = true
		}
	}

	if !hasUpper {
		return ErrorPasswordMissingUppercase
	}
	if !hasLower {
		return ErrorPasswordMissingLowercase
	}
	if !hasNumber {
		return ErrorPasswordMissingNumber
	}
	if !hasSpecialChar {
		return ErrorPasswordMissingSpecialChar
	}

	return nil
}
