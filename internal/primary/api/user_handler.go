package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mgwinsor/weekbyweek/internal/app/user"
)

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createUserResponse, err := s.userService.CreateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrEmailExists) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createUserResponse)
}
