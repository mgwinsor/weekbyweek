package user

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"dob"`
}

type CreateUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	DateOfBirth time.Time `json:"dob"`
}
