package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/mgwinsor/weekbyweek/internal/app"
	"github.com/mgwinsor/weekbyweek/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	dbURL := "file:data/dev.db?_fk=1"
	dbConn, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	dbQueries := database.New(dbConn)

	templates := template.Must(template.ParseGlob("web/templates/*.html"))

	svr := app.NewServer(dbQueries, templates)

	mux := http.NewServeMux()

	svr.SetupRoutes(mux)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
