package user

import (
	"context"
	"errors"

	"github.com/mgwinsor/weekbyweek/internal/domain/user"
)

var ErrEmailExists = errors.New("email already exists")

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
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

	resp := &UserResponse{
		ID:          newUser.ID(),
		Email:       newUser.Email(),
		Username:    newUser.Username(),
		DateOfBirth: newUser.DateOfBirth(),
	}

	return resp, nil
}
