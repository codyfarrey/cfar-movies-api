package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"cfar-movies-api/db"
	"cfar-movies-api/handlers"
)

var mdb *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from OS environment")
	}

	mdb, err = db.Connect()

	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer mdb.Close()

	movieHandler := handlers.NewMovieHandler(mdb)

	http.HandleFunc("GET /health", handlers.HealthHandler)
	http.HandleFunc("GET /movies", movieHandler.MoviesHandler)
	http.HandleFunc("GET /movie", movieHandler.GetRandomMovieHandler)
	http.HandleFunc("POST /movie", movieHandler.AddMovie)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
