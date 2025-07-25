package app

import (
	"html/template"
	"net/http"

	"github.com/mgwinsor/weekbyweek/internal/database"
)

type Server struct {
	db        database.Querier
	templates *template.Template
}

func NewServer(db database.Querier, templates *template.Template) *Server {
	return &Server{
		db:        db,
		templates: templates,
	}
}

func (s *Server) SetupRoutes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", s.index)

	mux.HandleFunc("GET /register", s.registerGet)
	mux.HandleFunc("POST /register", s.registerPost)

	mux.HandleFunc("GET /login", s.loginGet)
}
