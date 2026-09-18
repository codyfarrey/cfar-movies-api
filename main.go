package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from OS environment")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Println(fmt.Sprintf("%s %s %s %s", host, user, password, dbname))

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	fmt.Println("Hello World!")

	db, err := sql.Open("postgres", connStr)
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

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "OK")
}
