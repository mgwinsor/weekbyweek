package app

import (
	"database/sql"
	"log"
	"net/http"

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

	data := PageData{
		PageTitle:  "Create Your Account",
		FormValues: make(map[string]string),
		Errors:     make(map[string]string),
	}

	for k, v := range r.PostForm {
		data.FormValues[k] = v[0]
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	dob := r.FormValue("dob")

	if err := validateUsernameSyntax(username); err != nil {
		data.Errors["username"] = err.Error()
	}

	if err := validateEmailSyntax(email); err != nil {
		data.Errors["email"] = err.Error()
	}

	if err := valdiatePassword(password); err != nil {
		data.Errors["password"] = err.Error()
	}

	dateOfBirth, err := validateDateOfBirth(dob)
	if err != nil {
		data.Errors["dob"] = err.Error()
	}

	if _, exists := data.Errors["username"]; !exists {
		_, err := s.db.GetUserByUsername(r.Context(), username)
		if err == nil {
			data.Errors["username"] = ErrorUsernameExists.Error()
		} else if err != sql.ErrNoRows {
			http.Error(w, "Error reading username from database", http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}

	if _, exists := data.Errors["email"]; !exists {
		_, err := s.db.GetUserByEmail(r.Context(), email)
		if err == nil {
			data.Errors["email"] = ErrorEmailExists.Error()
		} else if err != sql.ErrNoRows {
			log.Print(err)
			http.Error(w, "Error reading email from database", http.StatusInternalServerError)
			return
		}
	}

	if len(data.Errors) > 0 {
		w.WriteHeader(http.StatusOK)
		s.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	passwordHash, err := s.authService.HashPassword(password)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	_, err = s.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		DateOfBirth:  dateOfBirth,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Redirect", "/home")
	w.WriteHeader(http.StatusCreated)
}
