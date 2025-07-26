package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/weekbyweek/internal/database"
)

func (s *Server) registerGet(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) registerPost(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := valdiatePassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dateOfBirth, err := time.Parse(time.DateOnly, r.FormValue("dob"))
	if err != nil {
		http.Error(w, "Invalid date of birth", http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUserByUsername(r.Context(), username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if user != (database.User{}) {
		http.Error(w, "Username already registered", http.StatusConflict)
		return
	}

	user, err = s.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if user != (database.User{}) {
		http.Error(w, "Email already registerd", http.StatusConflict)
		return
	}

	_, err = s.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: password,
		DateOfBirth:  dateOfBirth,
	})
	w.WriteHeader(http.StatusCreated)
}
