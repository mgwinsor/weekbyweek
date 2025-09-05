package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/app/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserAPIIntegration(t *testing.T) {
	validDOB := time.Date(1992, time.November, 21, 0, 0, 0, 0, time.UTC)
	id, _ := uuid.Parse("4762e4fb-b6bd-487d-834d-7a8c20c78be9")

	requestDTO := user.CreateUserRequest{
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: validDOB,
	}
	requestBody, _ := json.Marshal(requestDTO)

	responseDTO := user.CreateUserResponse{
		ID:          id,
		Email:       "john@example.com",
		Username:    "johndoe",
		DateOfBirth: validDOB,
	}
	responseBody, _ := json.Marshal(responseDTO)

	tests := []struct {
		name               string
		method             string
		path               string
		body               []byte
		mockSetup          func(m *MockUserService)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "successfully call users endpoint",
			method: http.MethodPost,
			path:   "/users",
			body:   requestBody,
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, requestDTO).
					Return(&responseDTO, nil).
					Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       string(responseBody),
		},
		{
			name:               "wrong method for users endpoint",
			method:             http.MethodGet,
			path:               "/users",
			body:               nil,
			mockSetup:          func(m *MockUserService) {},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "",
		},
		{
			name:               "call non-existent route",
			method:             http.MethodPost,
			path:               "/nonexistent",
			body:               nil,
			mockSetup:          func(m *MockUserService) {},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found",
		},
		{
			name:   "user service returns error",
			method: http.MethodPost,
			path:   "/users",
			body:   requestBody,
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, requestDTO).
					Return(nil, user.ErrEmailExists).
					Once()
			},
			expectedStatusCode: http.StatusConflict,
			expectedBody:       user.ErrEmailExists.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			server := NewUserHandler(mockService)
			router := server.RegisterRoutes()

			ts := httptest.NewServer(router)
			defer ts.Close()

			req, err := http.NewRequest(tt.method, ts.URL+tt.path, bytes.NewBuffer(tt.body))
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(string(respBodyBytes)))

			mockService.AssertExpectations(t)
		})
	}
}
