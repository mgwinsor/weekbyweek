package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type BcryptHasher struct{}

func (h BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func HashPassword(password string, hasher PasswordHasher) (string, error) {
	hash, err := hasher.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
