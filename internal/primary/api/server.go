package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mgwinsor/weekbyweek/internal/app/user"
)

type Server struct {
	userService user.Service
}

func NewServer(service user.Service) *Server {
	return &Server{
		userService: service,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", s.handleCreateUser)
	})

	return r
}
