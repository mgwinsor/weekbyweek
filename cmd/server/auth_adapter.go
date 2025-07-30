package main

import "github.com/mgwinsor/weekbyweek/internal/auth"

type authAdapter struct {
	hasher auth.PasswordHasher
}

func (a *authAdapter) HashPassword(password string) (string, error) {
	return auth.HashPassword(password, a.hasher)
}
