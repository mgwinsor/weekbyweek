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

type UserResponse struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DateOfBirth time.Time
}
