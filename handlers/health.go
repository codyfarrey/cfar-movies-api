package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HealthHandler(rw http.ResponseWriter, req *http.Request) {
	log.Print("Checking application health...")
	fmt.Fprintln(rw, "OK")
}
