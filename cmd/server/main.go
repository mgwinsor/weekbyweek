package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mgwinsor/weekbyweek/internal/app/user"
	"github.com/mgwinsor/weekbyweek/internal/primary/api"
	"github.com/mgwinsor/weekbyweek/internal/secondary/storage/memory"
)

func main() {
	userRepo := memory.NewUserRepository()
	userService := user.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userHandler.RegisterRoutes(r)

	log.Println("Server starting on port 8080")
	http.ListenAndServe(":8080", r)
}
