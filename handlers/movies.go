package handlers

import (
	"cfar-movies-api/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type MovieHandler struct {
	DB *sql.DB
}

func NewMovieHandler(db *sql.DB) *MovieHandler {
	return &MovieHandler{DB: db}
}

func (h *MovieHandler) MoviesHandler(rw http.ResponseWriter, req *http.Request) {
	log.Print("Returning all movies from database...")
	rows, err := h.DB.Query("SELECT tmdb_id, title, release_date, genres FROM movies")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
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

func (h *MovieHandler) AddMovie(rw http.ResponseWriter, req *http.Request) {
	var m models.Movie
	err := json.NewDecoder(req.Body).Decode(&m)

	log.Printf("Attempting to add movie %s to database", m.Title)

	if err != nil {
		http.Error(rw, "invalid request", http.StatusBadRequest)
		log.Print("Request is invalid: ", err)
		return
	}

	_, err = h.DB.Exec("INSERT INTO movies (tmdb_id, title, release_date, genres) VALUES ($1, $2, $3, $4)", m.ID, m.Title, m.ReleaseDate, m.Genres)

	if err != nil {
		http.Error(rw, "error writing to database", http.StatusInternalServerError)
		log.Print("error writing to database: ", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(m)

}

func (h *MovieHandler) GetRandomMovieHandler(rw http.ResponseWriter, req *http.Request) {
	row := h.DB.QueryRow("SELECT tmdb_id, title, release_date, genres FROM movies ORDER BY RANDOM() LIMIT 1")

	var m models.Movie
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
