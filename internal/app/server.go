package app

import (
	"html/template"
	"net/http"

	"github.com/mgwinsor/weekbyweek/internal/database"
)

type AuthService interface {
	HashPassword(password string) (string, error)
}

type PageData struct {
	PageTitle  string
	User       database.User
	FormValues map[string]string
	Errors     map[string]string
}

type Server struct {
	db          database.Querier
	authService AuthService
	templates   *template.Template
}

func NewServer(db database.Querier, authService AuthService, templates *template.Template) *Server {
	return &Server{
		db:          db,
		authService: authService,
		templates:   templates,
	}
}

func (s *Server) SetupRoutes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", s.index)

	mux.HandleFunc("GET /register", s.registerGet)
	mux.HandleFunc("POST /register", s.registerPost)

	mux.HandleFunc("GET /login", s.loginGet)

	mux.HandleFunc("GET /home", s.homeGet)
}
