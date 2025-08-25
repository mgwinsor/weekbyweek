package user

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email       string
	Username    string
	DateOfBirth time.Time
}

type CreateUserResponse struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DateOfBirth time.Time
}
