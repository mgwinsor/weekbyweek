package integration

import (
	"context"
	"testing"
	"time"

	"github.com/mgwinsor/weekbyweek/internal/app/user"
	"github.com/mgwinsor/weekbyweek/internal/secondary/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserIntegration(t *testing.T) {
	dob := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)

	newUserRequest := user.CreateUserRequest{
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: dob,
	}

	tests := []struct {
		name             string
		request          user.CreateUserRequest
		preExistingUsers []user.CreateUserRequest
		setupFn          func(*user.Service)
		wantErr          bool
	}{
		{
			name:             "successfully create user",
			request:          newUserRequest,
			preExistingUsers: nil,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := memory.NewUserRepository()
			userService := user.NewUserService(userRepo)

			for _, req := range tt.preExistingUsers {
				_, err := userService.CreateUser(context.Background(), req)
				require.NoError(t, err, "Setup failed: could not create pre-existing user")
			}

			response, err := userService.CreateUser(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)

				assert.Equal(t, tt.request.Email, response.Email)
				assert.Equal(t, tt.request.Username, response.Username)
				assert.Equal(t, tt.request.DateOfBirth, response.DateOfBirth)
				assert.NotEmpty(t, response.ID)

				savedUser, err := userRepo.FindByEmail(context.Background(), tt.request.Email)
				require.NoError(t, err, "User should be retrievable after creation")
				require.NotNil(t, savedUser)

				assert.Equal(t, tt.request.Email, savedUser.Email())
				assert.Equal(t, tt.request.Username, savedUser.Username())
				assert.Equal(t, tt.request.DateOfBirth, savedUser.DateOfBirth())
				assert.NotEmpty(t, savedUser.CreatedAt())
				assert.NotEmpty(t, savedUser.UpdatedAt())
			}
		})
	}
}
