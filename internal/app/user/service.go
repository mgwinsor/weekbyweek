package user

import (
	"context"
	"errors"

	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

var ErrEmailExists = errors.New("email already exists")

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
}

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(repo user.UserRepository) *userService {
	return &userService{userRepo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	newUser, err := user.NewUser(req.Email, req.Username, req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	resp := &CreateUserResponse{
		ID:          newUser.ID(),
		Email:       newUser.Email(),
		Username:    newUser.Username(),
		DateOfBirth: newUser.DateOfBirth(),
	}

	return resp, nil
}
