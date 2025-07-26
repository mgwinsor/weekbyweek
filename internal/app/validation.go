package app

import (
	"errors"
	"unicode"
)

var (
	ErrorPasswordTooShort           = errors.New("Password must be at least 12 characters long")
	ErrorPasswordMissingUppercase   = errors.New("Password must contain at least 1 uppercase letter")
	ErrorPasswordMissingLowercase   = errors.New("Password must contain at least 1 lowercase letter")
	ErrorPasswordMissingNumber      = errors.New("Password must contain at least 1 number")
	ErrorPasswordMissingSpecialChar = errors.New("Password must contain at least 1 special character")
)

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
