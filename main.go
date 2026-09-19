package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Genres      string `json:"genres"`
}

var db *sql.DB
var err error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from OS environment")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Printf("Connecting to database: [%s] as user: [%s]\n", dbname, user)

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dbErr := db.Ping()

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	log.Println("Connected to database!")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/movies", moviesHandler)
	http.HandleFunc("/movie", getRandomMovieHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "OK")
}

func moviesHandler(rw http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT tmdb_id, title, release_date, genres FROM movies")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseDate, &m.Genres)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		movies = append(movies, m)
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(movies)
}

func getRandomMovieHandler(rw http.ResponseWriter, req *http.Request) {
	row := db.QueryRow("SELECT tmdb_id, title, release_date, genres FROM movies ORDER BY RANDOM() LIMIT 1")

	var m Movie
	err := row.Scan(&m.ID, &m.Title, &m.ReleaseDate, &m.Genres)
	if err == sql.ErrNoRows {
		http.Error(rw, "no movies found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(m)
}
