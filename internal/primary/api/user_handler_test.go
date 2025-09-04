package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/app/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req user.CreateUserRequest) (*user.CreateUserResponse, error) {
	args := m.Called(ctx, req)

	var resp *user.CreateUserResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*user.CreateUserResponse)
	}

	return resp, args.Error(1)
}

func TestHandleCreateUser(t *testing.T) {
	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	id, _ := uuid.Parse("4762e4fb-b6bd-487d-834d-7a8c20c78be9")

	requestDTO := user.CreateUserRequest{
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: validDOB,
	}
	requestBody, _ := json.Marshal(requestDTO)

	successResponseDTO := user.CreateUserResponse{
		ID:          id,
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: validDOB,
	}

	successResponseBody, _ := json.Marshal(successResponseDTO)

	tests := []struct {
		name               string
		inputBody          string
		mockSetup          func(MockService *MockUserService)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "successfully created user",
			inputBody: string(requestBody),
			mockSetup: func(MockService *MockUserService) {
				MockService.On("CreateUser", mock.Anything, requestDTO).
					Return(&successResponseDTO, nil).
					Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       string(successResponseBody),
		},
		{
			name:      "email exists error",
			inputBody: string(requestBody),
			mockSetup: func(MockService *MockUserService) {
				MockService.On("CreateUser", mock.Anything, requestDTO).
					Return(nil, user.ErrEmailExists).
					Once()
			},
			expectedStatusCode: http.StatusConflict,
			expectedBody:       "email already exists",
		},
		{
			name:               "email is not a string",
			inputBody:          newCreateUserPayload(map[string]any{"email": 1234}),
			mockSetup:          func(MockService *MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Invalid request body",
		},
		{
			name:               "invalid date of birth format",
			inputBody:          newCreateUserPayload(map[string]any{"dob": "21-11-1992"}),
			mockSetup:          func(MockService *MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Invalid request body",
		},
		{
			name:      "unexpected error",
			inputBody: string(requestBody),
			mockSetup: func(MockService *MockUserService) {
				MockService.On("CreateUser", mock.Anything, requestDTO).
					Return(nil, errors.New("unexpected error")).
					Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "Failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			server := NewServer(mockService)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tt.inputBody))

			server.handleCreateUser(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code, "status code should match expected")
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()), "response body should match expected")

			mockService.AssertExpectations(t)
		})
	}
}

func newCreateUserPayload(overrides map[string]any) string {
	payload := map[string]any{
		"email":    "test@example.com",
		"username": "testuser",
		"dob":      "1992-11-21T00:00:00Z",
	}

	maps.Copy(payload, overrides)

	body, _ := json.Marshal(payload)
	return string(body)
}
