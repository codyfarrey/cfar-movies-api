package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/codyfarrey/cfar-movies-api/db"
	"github.com/codyfarrey/cfar-movies-api/handlers"
	"github.com/codyfarrey/cfar-movies-api/security"
)

var mdb *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from OS environment")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}

	mdb, err = db.Connect()

	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer mdb.Close()

	movieHandler := handlers.NewMovieHandler(mdb)
	sp := security.NewSecurityPayload()

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/movies", sp.ValidateApiKey(movieHandler.MoviesHandler))
	http.HandleFunc("/movie", sp.ValidateApiKey(movieHandler.GetRandomMovieHandler))
	http.HandleFunc("POST /movie", sp.ValidateApiKey(movieHandler.AddMovie))

	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
