package app

import (
	"log"
	"net/http"
)

func (s *Server) loginGet(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
